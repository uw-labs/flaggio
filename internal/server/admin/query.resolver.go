package admin

import (
	"context"

	"github.com/victorkohl/flaggio/internal/flaggio"
)

var _ QueryResolver = &queryResolver{}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Ping(ctx context.Context) (*bool, error) {
	pong := true
	return &pong, nil
}

func (r *queryResolver) Flags(ctx context.Context, offset, limit *int) ([]*flaggio.Flag, error) {
	if limit == nil {
		v := 50
		limit = &v
	}
	var ofst, lmt *int64
	if offset != nil {
		v := int64(*offset)
		ofst = &v
	}
	if limit != nil {
		v := int64(*limit)
		lmt = &v
	}
	return r.flagRepo.FindAll(ctx, ofst, lmt)
}

func (r *queryResolver) Flag(ctx context.Context, id string) (*flaggio.Flag, error) {
	return r.flagRepo.FindByID(ctx, id)
}

func (r *queryResolver) Segments(ctx context.Context, offset, limit *int) ([]*flaggio.Segment, error) {
	if limit == nil {
		v := 50
		limit = &v
	}
	var ofst, lmt *int64
	if offset != nil {
		v := int64(*offset)
		ofst = &v
	}
	if limit != nil {
		v := int64(*limit)
		lmt = &v
	}
	return r.segmentRepo.FindAll(ctx, ofst, lmt)
}

func (r *queryResolver) Segment(ctx context.Context, id string) (*flaggio.Segment, error) {
	return r.segmentRepo.FindByID(ctx, id)
}
