package repository

import (
	"context"

	"github.com/victorkohl/flaggio/internal/flaggio"
)

type Variant interface {
	Create(ctx context.Context, flagID string, v flaggio.NewVariant) (string, error)
	Update(ctx context.Context, flagID, id string, v flaggio.UpdateVariant) error
	Delete(ctx context.Context, flagID, id string) error
}
