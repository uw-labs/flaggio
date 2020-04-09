package repository

//go:generate mockgen -destination=./mocks/rule_mock.go -package=repository_mock github.com/victorkt/flaggio/internal/repository Rule

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

// Rule represents a set of operations available to list and manage rules.
type Rule interface {
	// FindFlagRuleByID returns a flag rule that has a given ID.
	FindFlagRuleByID(ctx context.Context, flagIDHex, idHex string) (*flaggio.FlagRule, error)
	// CreateFlagRule creates a new rule under a flag.
	CreateFlagRule(ctx context.Context, flagID string, input flaggio.NewFlagRule) (string, error)
	// UpdateFlagRule updates a rule under a flag.
	UpdateFlagRule(ctx context.Context, flagID, id string, input flaggio.UpdateFlagRule) error
	// DeleteFlagRule deletes a rule under a flag.
	DeleteFlagRule(ctx context.Context, flagID, id string) error
	// FindSegmentRuleByID returns a segment rule that has a given ID.
	FindSegmentRuleByID(ctx context.Context, segmentIDHex, idHex string) (*flaggio.SegmentRule, error)
	// CreateSegmentRule creates a new rule under a segment.
	CreateSegmentRule(ctx context.Context, segmentID string, input flaggio.NewSegmentRule) (string, error)
	// UpdateSegmentRule updates a rule under a segment.
	UpdateSegmentRule(ctx context.Context, segmentID, id string, input flaggio.UpdateSegmentRule) error
	// DeleteSegmentRule deletes a rule under a segment.
	DeleteSegmentRule(ctx context.Context, segmentID, id string) error
}
