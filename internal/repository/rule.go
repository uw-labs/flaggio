package repository

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

// Rule represents a set of operations available to list and manage rules.
type Rule interface {
	// CreateFlagRule creates a new rule under a flag.
	CreateFlagRule(ctx context.Context, flagID string, input flaggio.NewFlagRule) (string, error)
	// UpdateFlagRule updates a rule under a flag.
	UpdateFlagRule(ctx context.Context, flagID, id string, input flaggio.UpdateFlagRule) error
	// DeleteFlagRule deletes a rule under a flag.
	DeleteFlagRule(ctx context.Context, flagID, id string) error
	// CreateSegmentRule creates a new rule under a segment.
	CreateSegmentRule(ctx context.Context, segmentID string, input flaggio.NewSegmentRule) (string, error)
	// UpdateSegmentRule updates a rule under a segment.
	UpdateSegmentRule(ctx context.Context, segmentID, id string, input flaggio.UpdateSegmentRule) error
	// DeleteSegmentRule deletes a rule under a segment.
	DeleteSegmentRule(ctx context.Context, segmentID, id string) error
}
