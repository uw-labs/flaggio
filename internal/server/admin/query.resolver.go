package admin

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

var _ QueryResolver = &queryResolver{}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Ping(ctx context.Context) (bool, error) {
	return true, nil
}

func (r *queryResolver) Flags(ctx context.Context, offset, limit *int) ([]*flaggio.Flag, error) {
	var ofst, lmt *int64
	if offset != nil {
		v := int64(*offset)
		ofst = &v
	}
	if limit != nil {
		v := int64(*limit)
		lmt = &v
	}
	return r.FlagRepo.FindAll(ctx, ofst, lmt)
}

func (r *queryResolver) Flag(ctx context.Context, id string) (*flaggio.Flag, error) {
	return r.FlagRepo.FindByID(ctx, id)
}

func (r *queryResolver) Segments(ctx context.Context, offset, limit *int) ([]*flaggio.Segment, error) {
	var ofst, lmt *int64
	if offset != nil {
		v := int64(*offset)
		ofst = &v
	}
	if limit != nil {
		v := int64(*limit)
		lmt = &v
	}
	return r.SegmentRepo.FindAll(ctx, ofst, lmt)
}

func (r *queryResolver) Segment(ctx context.Context, id string) (*flaggio.Segment, error) {
	return r.SegmentRepo.FindByID(ctx, id)
}
