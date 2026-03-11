package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"time"
)

// Close releases the shared connection and cleans up per-engine scratch space.
func (db *Database) Close() error {
	if db == nil {
		return nil
	}

	var firstErr error
	if db.DB != nil {
		if err := db.DB.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if db.duckCleanup != nil {
		if err := db.duckCleanup(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

type serializedJob struct {
	ctx    context.Context
	fn     func(context.Context, *sql.DB) (any, error)
	result chan serializedResult
	kind   WorkloadKind
}

type serializedResult struct {
	value any
	err   error
}

type WorkloadKind int

const (
	WorkloadGeneral WorkloadKind = iota
	WorkloadWebRead
	WorkloadUserUpload
	WorkloadArchive
	WorkloadRealtime
)

type serializedPipeline struct {
	lanes []chan serializedJob
	order []int
	turn  int
}

type duckDBMaintenance struct {
	db   *Database
	jobs chan duckDBMaintenanceJob
}

type duckDBMaintenanceJob struct {
	ctx  context.Context
	logf func(string, ...any)
	done chan error
}

const serializedWaitFloor = 45 * time.Second

func startSerializedPipeline(db *sql.DB) *serializedPipeline {
	p := &serializedPipeline{
		lanes: []chan serializedJob{
			make(chan serializedJob, 64),
			make(chan serializedJob, 128),
			make(chan serializedJob, 32),
			make(chan serializedJob, 32),
			make(chan serializedJob, 128),
		},
		order: []int{int(WorkloadRealtime), int(WorkloadArchive), int(WorkloadUserUpload), int(WorkloadWebRead), int(WorkloadGeneral)},
	}

	getNextJob := func() serializedJob {
		for i := 0; i < len(p.order); i++ {
			laneIdx := p.order[p.turn%len(p.order)]
			p.turn++
			select {
			case job := <-p.lanes[laneIdx]:
				return job
			default:
			}
		}

		select {
		case job := <-p.lanes[p.order[0]]:
			return job
		case job := <-p.lanes[p.order[1]]:
			return job
		case job := <-p.lanes[p.order[2]]:
			return job
		case job := <-p.lanes[p.order[3]]:
			return job
		case job := <-p.lanes[p.order[4]]:
			return job
		}
	}

	go func() {
		for {
			job := getNextJob()
			select {
			case <-job.ctx.Done():
				job.result <- serializedResult{nil, job.ctx.Err()}
				continue
			default:
			}

			value, err := job.fn(job.ctx, db)
			job.result <- serializedResult{value, err}
		}
	}()

	return p
}

func startDuckDBMaintenance(db *Database) *duckDBMaintenance {
	m := &duckDBMaintenance{db: db, jobs: make(chan duckDBMaintenanceJob, 4)}

	go func() {
		for job := range m.jobs {
			logf := job.logf
			if logf == nil {
				logf = log.Printf
			}
			if job.ctx == nil {
				job.ctx = context.Background()
			}

			select {
			case <-job.ctx.Done():
				job.done <- job.ctx.Err()
				close(job.done)
				continue
			default:
			}

			maintenanceCtx, cancel := queueFriendlyContext(job.ctx, 90*time.Minute)
			err := m.db.withSerializedConnectionFor(maintenanceCtx, WorkloadArchive, func(runCtx context.Context, conn *sql.DB) error {
				return runDuckDBMaintenance(runCtx, conn, logf)
			})
			cancel()

			select {
			case job.done <- err:
			default:
			}
			close(job.done)
		}
	}()

	return m
}

func (m *duckDBMaintenance) enqueue(ctx context.Context, logf func(string, ...any)) <-chan error {
	if m == nil {
		done := make(chan error, 1)
		close(done)
		return done
	}
	job := duckDBMaintenanceJob{ctx: ctx, logf: logf, done: make(chan error, 1)}
	for {
		select {
		case <-ctx.Done():
			job.done <- ctx.Err()
			close(job.done)
			return job.done
		case m.jobs <- job:
			return job.done
		default:
			runtime.Gosched()
		}
	}
}

func (db *Database) ScheduleDuckDBMaintenance(ctx context.Context, logf func(string, ...any)) <-chan error {
	if db == nil || db.Driver != "duckdb" || db.upkeep == nil {
		done := make(chan error, 1)
		close(done)
		return done
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return db.upkeep.enqueue(ctx, logf)
}

func (db *Database) serializedEnabled() bool {
	if db == nil {
		return false
	}
	switch db.Driver {
	case "duckdb", "sqlite", "chai":
		return true
	default:
		return false
	}
}

func (p *serializedPipeline) enqueue(ctx context.Context, job serializedJob) error {
	lane := p.laneFor(job.kind)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case lane <- job:
			return nil
		default:
			runtime.Gosched()
		}
	}
}

func (p *serializedPipeline) laneFor(kind WorkloadKind) chan serializedJob {
	if int(kind) >= 0 && int(kind) < len(p.lanes) {
		return p.lanes[kind]
	}
	return p.lanes[WorkloadGeneral]
}

func queueFriendlyContext(ctx context.Context, min time.Duration) (context.Context, context.CancelFunc) {
	if min <= 0 {
		min = serializedWaitFloor
	}
	if ctx == nil {
		return context.WithTimeout(context.Background(), min)
	}
	if deadline, ok := ctx.Deadline(); ok {
		if time.Until(deadline) >= min {
			return ctx, func() {}
		}
		return context.WithTimeout(context.WithoutCancel(ctx), min)
	}
	return context.WithTimeout(ctx, min)
}

func (db *Database) withSerializedConnection(ctx context.Context, fn func(context.Context, *sql.DB) error) error {
	return db.withSerializedConnectionFor(ctx, WorkloadGeneral, fn)
}

func (db *Database) withSerializedConnectionFor(ctx context.Context, kind WorkloadKind, fn func(context.Context, *sql.DB) error) error {
	if db == nil || db.DB == nil {
		return fmt.Errorf("database unavailable")
	}

	if !db.serializedEnabled() || db.pipeline == nil {
		return fn(ctx, db.DB)
	}

	job := serializedJob{
		ctx:    ctx,
		fn:     func(c context.Context, conn *sql.DB) (any, error) { return nil, fn(c, conn) },
		result: make(chan serializedResult, 1),
		kind:   kind,
	}

	if err := db.pipeline.enqueue(ctx, job); err != nil {
		return err
	}

	select {
	case res := <-job.result:
		return res.err
	case <-ctx.Done():
		return ctx.Err()
	}
}
