// services_web.go defines service interfaces and a database-backed implementation for web handlers.
package httpapi

import (
	"context"
	"fmt"
	"strings"

	"safecast-new-map/pkg/database"
)

// MarkerService encapsulates marker-related API use-cases (markers with spectra, update coordinates).
type MarkerService interface {
	GetMarkersWithSpectra(ctx context.Context, bounds database.Bounds) ([]database.Marker, error)
	UpdateTrackCoordinates(ctx context.Context, trackID string, lat, lon float64) (int64, error)
}

// TrackService encapsulates track metadata and bounds (track info, bounding box for tracks).
type TrackService interface {
	GetTrackInfo(ctx context.Context, trackID string) (database.TrackInfo, error)
	GetTrackBounds(ctx context.Context, trackID string) (database.Bounds, bool, error)
}

// Services groups the marker and track services used by web handlers.
type Services struct {
	Marker MarkerService
	Track  TrackService
}

type databaseService struct {
	db *database.Database
}

func newDatabaseServices(db *database.Database) Services {
	svc := &databaseService{db: db}
	return Services{
		Marker: svc,
		Track:  svc,
	}
}

func (s *databaseService) GetMarkersWithSpectra(ctx context.Context, bounds database.Bounds) ([]database.Marker, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("database not available")
	}
	return s.db.GetMarkersWithSpectra(ctx, bounds)
}

func (s *databaseService) UpdateTrackCoordinates(ctx context.Context, trackID string, lat, lon float64) (int64, error) {
	if s == nil || s.db == nil {
		return 0, fmt.Errorf("database not available")
	}
	return s.db.UpdateTrackCoordinates(ctx, trackID, lat, lon)
}

func (s *databaseService) GetTrackInfo(ctx context.Context, trackID string) (database.TrackInfo, error) {
	if s == nil || s.db == nil {
		return database.TrackInfo{}, fmt.Errorf("database not available")
	}
	return s.db.GetTrackInfo(ctx, trackID)
}

func (s *databaseService) GetTrackBounds(ctx context.Context, trackID string) (database.Bounds, bool, error) {
	if s == nil || s.db == nil {
		return database.Bounds{}, false, fmt.Errorf("database not available")
	}
	return s.db.GetTrackBounds(ctx, trackID)
}

// normalizeTrackIDs trims and drops empty strings from a list of track IDs.
func normalizeTrackIDs(raw []string) []string {
	out := make([]string, 0, len(raw))
	for _, id := range raw {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
