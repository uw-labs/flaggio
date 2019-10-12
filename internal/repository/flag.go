package repository

import (
	"context"

	"github.com/victorkohl/flaggio/internal/flaggio"
)

type Flag interface {
	FindAll(ctx context.Context, offset, limit *int64) ([]*flaggio.Flag, error)
	FindByID(ctx context.Context, id string) (*flaggio.Flag, error)
	FindByKey(ctx context.Context, key string) (*flaggio.Flag, error)
	Create(ctx context.Context, f flaggio.NewFlag) (*flaggio.Flag, error)
	// Update(ctx context.Context, id string, f flaggio.UpdateFlag) (*flaggio.Flag, error)
	Delete(ctx context.Context, id string) error
}
