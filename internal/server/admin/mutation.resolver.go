package admin

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

var _ MutationResolver = &mutationResolver{}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Ping(ctx context.Context) (bool, error) {
	return true, nil
}

func (r *mutationResolver) CreateFlag(ctx context.Context, input flaggio.NewFlag) (*flaggio.Flag, error) {
	id, err := r.FlagRepo.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return r.FlagRepo.FindByID(ctx, id)
}

func (r *mutationResolver) UpdateFlag(ctx context.Context, id string, input flaggio.UpdateFlag) (*flaggio.Flag, error) {
	if err := r.FlagRepo.Update(ctx, id, input); err != nil {
		return nil, err
	}
	return r.FlagRepo.FindByID(ctx, id)
}

func (r *mutationResolver) DeleteFlag(ctx context.Context, id string) (string, error) {
	err := r.FlagRepo.Delete(ctx, id)
	return id, err
}

func (r *mutationResolver) CreateVariant(ctx context.Context, flagID string, input flaggio.NewVariant) (*flaggio.Variant, error) {
	return r.VariantRepo.Create(ctx, flagID, input)
}

func (r *mutationResolver) UpdateVariant(ctx context.Context, flagID, id string, input flaggio.UpdateVariant) (*flaggio.Variant, error) {
	return r.VariantRepo.Update(ctx, flagID, id, input)
}

func (r *mutationResolver) DeleteVariant(ctx context.Context, flagID, id string) (string, error) {
	err := r.VariantRepo.Delete(ctx, flagID, id)
	return id, err
}

func (r *mutationResolver) CreateFlagRule(ctx context.Context, flagID string, input flaggio.NewFlagRule) (string, error) {
	return r.RuleRepo.CreateFlagRule(ctx, flagID, input)
}

func (r *mutationResolver) UpdateFlagRule(ctx context.Context, flagID, id string, input flaggio.UpdateFlagRule) (bool, error) {
	err := r.RuleRepo.UpdateFlagRule(ctx, flagID, id, input)
	return err == nil, err
}

func (r *mutationResolver) DeleteFlagRule(ctx context.Context, flagID, id string) (string, error) {
	err := r.RuleRepo.DeleteFlagRule(ctx, flagID, id)
	return id, err
}

func (r *mutationResolver) CreateSegmentRule(ctx context.Context, segmentID string, input flaggio.NewSegmentRule) (string, error) {
	return r.RuleRepo.CreateSegmentRule(ctx, segmentID, input)
}

func (r *mutationResolver) UpdateSegmentRule(ctx context.Context, segmentID, id string, input flaggio.UpdateSegmentRule) (bool, error) {
	err := r.RuleRepo.UpdateSegmentRule(ctx, segmentID, id, input)
	return err == nil, err
}

func (r *mutationResolver) DeleteSegmentRule(ctx context.Context, segmentID, id string) (string, error) {
	err := r.RuleRepo.DeleteSegmentRule(ctx, segmentID, id)
	return id, err
}

func (r *mutationResolver) CreateSegment(ctx context.Context, input flaggio.NewSegment) (*flaggio.Segment, error) {
	return r.SegmentRepo.Create(ctx, input)
}

func (r *mutationResolver) UpdateSegment(ctx context.Context, id string, input flaggio.UpdateSegment) (bool, error) {
	err := r.SegmentRepo.Update(ctx, id, input)
	return err == nil, err
}

func (r *mutationResolver) DeleteSegment(ctx context.Context, id string) (string, error) {
	err := r.SegmentRepo.Delete(ctx, id)
	return id, err
}
