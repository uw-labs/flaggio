package admin

import (
	"context"

	"github.com/victorkohl/flaggio/internal/flaggio"
)

var _ MutationResolver = &mutationResolver{}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Ping(ctx context.Context) (*bool, error) {
	pong := true
	return &pong, nil
}

func (r *mutationResolver) CreateFlag(ctx context.Context, input flaggio.NewFlag) (*flaggio.Flag, error) {
	return r.flagRepo.Create(ctx, input)
}

func (r *mutationResolver) UpdateFlag(ctx context.Context, id string, input flaggio.UpdateFlag) (*flaggio.Flag, error) {
	return r.flagRepo.Update(ctx, id, input)
}

func (r *mutationResolver) DeleteFlag(ctx context.Context, id string) (string, error) {
	err := r.flagRepo.Delete(ctx, id)
	return id, err
}

func (r *mutationResolver) CreateVariant(ctx context.Context, flagID string, input flaggio.NewVariant) (string, error) {
	return r.variantRepo.Create(ctx, flagID, input)
}

func (r *mutationResolver) UpdateVariant(ctx context.Context, flagID, id string, input flaggio.UpdateVariant) (bool, error) {
	err := r.variantRepo.Update(ctx, flagID, id, input)
	return err == nil, err
}

func (r *mutationResolver) DeleteVariant(ctx context.Context, flagID, id string) (string, error) {
	err := r.variantRepo.Delete(ctx, flagID, id)
	return id, err
}

func (r *mutationResolver) CreateFlagRule(ctx context.Context, flagID string, input flaggio.NewFlagRule) (string, error) {
	return r.ruleRepo.CreateFlagRule(ctx, flagID, input)
}

func (r *mutationResolver) UpdateFlagRule(ctx context.Context, flagID, id string, input flaggio.UpdateFlagRule) (bool, error) {
	err := r.ruleRepo.UpdateFlagRule(ctx, flagID, id, input)
	return err == nil, err
}

func (r *mutationResolver) DeleteFlagRule(ctx context.Context, flagID, id string) (string, error) {
	err := r.ruleRepo.DeleteFlagRule(ctx, flagID, id)
	return id, err
}

func (r *mutationResolver) CreateSegmentRule(ctx context.Context, segmentID string, input flaggio.NewSegmentRule) (string, error) {
	return r.ruleRepo.CreateSegmentRule(ctx, segmentID, input)
}

func (r *mutationResolver) UpdateSegmentRule(ctx context.Context, segmentID, id string, input flaggio.UpdateSegmentRule) (bool, error) {
	err := r.ruleRepo.UpdateSegmentRule(ctx, segmentID, id, input)
	return err == nil, err
}

func (r *mutationResolver) DeleteSegmentRule(ctx context.Context, segmentID, id string) (string, error) {
	err := r.ruleRepo.DeleteSegmentRule(ctx, segmentID, id)
	return id, err
}

func (r *mutationResolver) CreateSegment(ctx context.Context, input flaggio.NewSegment) (*flaggio.Segment, error) {
	return r.segmentRepo.Create(ctx, input)
}

func (r *mutationResolver) UpdateSegment(ctx context.Context, id string, input flaggio.UpdateSegment) (bool, error) {
	err := r.segmentRepo.Update(ctx, id, input)
	return err == nil, err
}

func (r *mutationResolver) DeleteSegment(ctx context.Context, id string) (string, error) {
	err := r.segmentRepo.Delete(ctx, id)
	return id, err
}
