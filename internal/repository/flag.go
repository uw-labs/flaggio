package repository

//go:generate mockgen -destination=./mocks/flag_mock.go -package=repository_mock github.com/uw-labs/flaggio/internal/repository Flag

import (
	"context"

	"github.com/uw-labs/flaggio/internal/flaggio"
)

// Flag represents a set of operations available to list and manage flags.
type Flag interface {
	// FindAll returns a list of flags, based on an optional offset and limit.
	FindAll(ctx context.Context, search *string, offset, limit *int64) (*flaggio.FlagResults, error)
	// FindByID returns a flag that has a given ID.
	FindByID(ctx context.Context, id string) (*flaggio.Flag, error)
	// FindByKey returns a flag that has a given key.
	FindByKey(ctx context.Context, key string) (*flaggio.Flag, error)
	// Create creates a new flag.
	Create(ctx context.Context, input flaggio.NewFlag) (string, error)
	// Update updates a flag.
	Update(ctx context.Context, id string, input flaggio.UpdateFlag) error
	// Delete deletes a flag.
	Delete(ctx context.Context, id string) error
}
