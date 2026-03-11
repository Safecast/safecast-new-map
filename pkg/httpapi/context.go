// context.go provides context helpers used by handlers and main (e.g. import shield).
package httpapi

import (
	"context"
	"time"
)

// WithMinimumDeadline returns a context that has at least min duration until its deadline.
// If ctx already has a deadline that is at least min away, ctx is returned unchanged with a no-op cancel.
// Otherwise a new context is created with the minimum duration so the DB pipeline has time to run.
func WithMinimumDeadline(ctx context.Context, min time.Duration) (context.Context, context.CancelFunc) {
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
