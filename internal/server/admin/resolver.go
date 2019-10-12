package admin

import (
	"context"

	"github.com/victorkohl/flaggio/internal/flaggio"
	"github.com/victorkohl/flaggio/internal/repository"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	flagRepo    repository.Flag
	variantRepo repository.Variant
}

func NewResolver(
	flagRepo repository.Flag,
	variantRepo repository.Variant,
) *Resolver {
	return &Resolver{
		flagRepo:    flagRepo,
		variantRepo: variantRepo,
	}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Ping(ctx context.Context) (*bool, error) {
	pong := true
	return &pong, nil
}

func (r *mutationResolver) CreateFlag(ctx context.Context, input flaggio.NewFlag) (*flaggio.Flag, error) {
	return r.flagRepo.Create(ctx, input)
}

func (r *mutationResolver) DeleteFlag(ctx context.Context, id string) (bool, error) {
	err := r.flagRepo.Delete(ctx, id)
	return err == nil, err
}

func (r *mutationResolver) CreateVariant(ctx context.Context, flagID string, input flaggio.NewVariant) (*flaggio.Variant, error) {
	return r.variantRepo.Create(ctx, flagID, input)
}

func (r *mutationResolver) DeleteVariant(ctx context.Context, flagID, id string) (bool, error) {
	err := r.variantRepo.Delete(ctx, flagID, id)
	return err == nil, err
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Ping(ctx context.Context) (*bool, error) {
	pong := true
	return &pong, nil
}

func (r *queryResolver) Flags(ctx context.Context) ([]*flaggio.Flag, error) {
	return r.flagRepo.FindAll(ctx, nil, nil)
}
