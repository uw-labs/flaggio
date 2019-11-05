package repository

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

type Segment interface {
	FindAll(ctx context.Context, offset, limit *int64) ([]*flaggio.Segment, error)
	FindByID(ctx context.Context, id string) (*flaggio.Segment, error)
	FindByKey(ctx context.Context, key string) (*flaggio.Segment, error)
	Create(ctx context.Context, f flaggio.NewSegment) (*flaggio.Segment, error)
	Update(ctx context.Context, id string, f flaggio.UpdateSegment) error
	Delete(ctx context.Context, id string) error
}
