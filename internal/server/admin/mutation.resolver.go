package admin

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

var _ MutationResolver = &mutationResolver{}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Ping(_ context.Context) (bool, error) {
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
	id, err := r.VariantRepo.Create(ctx, flagID, input)
	if err != nil {
		return nil, err
	}
	return r.VariantRepo.FindByID(ctx, flagID, id)
}

func (r *mutationResolver) UpdateVariant(ctx context.Context, flagID, id string, input flaggio.UpdateVariant) (*flaggio.Variant, error) {
	if err := r.VariantRepo.Update(ctx, flagID, id, input); err != nil {
		return nil, err
	}
	return r.VariantRepo.FindByID(ctx, flagID, id)
}

func (r *mutationResolver) DeleteVariant(ctx context.Context, flagID, id string) (string, error) {
	err := r.VariantRepo.Delete(ctx, flagID, id)
	return id, err
}

func (r *mutationResolver) CreateFlagRule(ctx context.Context, flagID string, input flaggio.NewFlagRule) (*flaggio.FlagRule, error) {
	id, err := r.RuleRepo.CreateFlagRule(ctx, flagID, input)
	if err != nil {
		return nil, err
	}
	return r.RuleRepo.FindFlagRuleByID(ctx, flagID, id)
}

func (r *mutationResolver) UpdateFlagRule(ctx context.Context, flagID, id string, input flaggio.UpdateFlagRule) (*flaggio.FlagRule, error) {
	if err := r.RuleRepo.UpdateFlagRule(ctx, flagID, id, input); err != nil {
		return nil, err
	}
	return r.RuleRepo.FindFlagRuleByID(ctx, flagID, id)
}

func (r *mutationResolver) DeleteFlagRule(ctx context.Context, flagID, id string) (string, error) {
	err := r.RuleRepo.DeleteFlagRule(ctx, flagID, id)
	return id, err
}

func (r *mutationResolver) CreateSegmentRule(ctx context.Context, segmentID string, input flaggio.NewSegmentRule) (*flaggio.SegmentRule, error) {
	id, err := r.RuleRepo.CreateSegmentRule(ctx, segmentID, input)
	if err != nil {
		return nil, err
	}
	return r.RuleRepo.FindSegmentRuleByID(ctx, segmentID, id)
}

func (r *mutationResolver) UpdateSegmentRule(ctx context.Context, segmentID, id string, input flaggio.UpdateSegmentRule) (*flaggio.SegmentRule, error) {
	if err := r.RuleRepo.UpdateSegmentRule(ctx, segmentID, id, input); err != nil {
		return nil, err
	}
	return r.RuleRepo.FindSegmentRuleByID(ctx, segmentID, id)
}

func (r *mutationResolver) DeleteSegmentRule(ctx context.Context, segmentID, id string) (string, error) {
	err := r.RuleRepo.DeleteSegmentRule(ctx, segmentID, id)
	return id, err
}

func (r *mutationResolver) CreateSegment(ctx context.Context, input flaggio.NewSegment) (*flaggio.Segment, error) {
	id, err := r.SegmentRepo.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return r.SegmentRepo.FindByID(ctx, id)
}

func (r *mutationResolver) UpdateSegment(ctx context.Context, id string, input flaggio.UpdateSegment) (*flaggio.Segment, error) {
	if err := r.SegmentRepo.Update(ctx, id, input); err != nil {
		return nil, err
	}
	return r.SegmentRepo.FindByID(ctx, id)
}

func (r *mutationResolver) DeleteSegment(ctx context.Context, id string) (string, error) {
	err := r.SegmentRepo.Delete(ctx, id)
	return id, err
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (string, error) {
	if err := r.UserRepo.Delete(ctx, id); err != nil {
		return id, err
	}
	if err := r.EvaluationRepo.DeleteAllByUserID(ctx, id); err != nil {
		return id, err
	}
	return id, nil
}

func (r *mutationResolver) DeleteEvaluation(ctx context.Context, id string) (string, error) {
	err := r.EvaluationRepo.DeleteByID(ctx, id)
	return id, err
}
