package repository

import (
	"context"

	"github.com/victorkohl/flaggio/internal/flaggio"
)

type Variant interface {
	Create(ctx context.Context, flagID string, v flaggio.NewVariant) (*flaggio.Variant, error)
	// Update(ctx context.Context, flagID string, v UpdateVariantDto) (*Variant, error)
	Delete(ctx context.Context, flagID, id string) error
}
