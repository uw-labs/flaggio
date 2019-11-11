package repository

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

// Segment represents a set of operations available to list and manage segments.
type Segment interface {
	// FindAll returns a list of segments, based on an optional offset and limit.
	FindAll(ctx context.Context, offset, limit *int64) ([]*flaggio.Segment, error)
	// FindByID returns a segment that has a given ID.
	FindByID(ctx context.Context, id string) (*flaggio.Segment, error)
	// Create creates a new segment.
	Create(ctx context.Context, input flaggio.NewSegment) (string, error)
	// Update updates a segment.
	Update(ctx context.Context, id string, input flaggio.UpdateSegment) error
	// Delete deletes a segment.
	Delete(ctx context.Context, id string) error
}
