package repository

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

// Variant represents a set of operations available to list and manage variants.
type Variant interface {
	// FindByID returns a variant that has a given ID.
	FindByID(ctx context.Context, flagIDHex, idHex string) (*flaggio.Variant, error)
	// Create creates a new variant under a flag.
	Create(ctx context.Context, flagID string, input flaggio.NewVariant) (string, error)
	// Update updates a variant under a flag.
	Update(ctx context.Context, flagID, id string, input flaggio.UpdateVariant) error
	// Delete deletes a variant under a flag.
	Delete(ctx context.Context, flagID, id string) error
}
