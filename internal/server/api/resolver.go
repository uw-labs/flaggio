package api

import (
	"context"

	"github.com/victorkohl/flaggio/internal/errors"
	"github.com/victorkohl/flaggio/internal/flaggio"
	"github.com/victorkohl/flaggio/internal/repository"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	flagRepo repository.Flag
}

func NewResolver(
	flagRepo repository.Flag,
) *Resolver {
	return &Resolver{
		flagRepo: flagRepo,
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

type queryResolver struct{ *Resolver }

func (r *queryResolver) Ping(ctx context.Context) (*bool, error) {
	pong := true
	return &pong, nil
}

func (r *queryResolver) Evaluate(ctx context.Context, flagKey string, userID string, context map[string]interface{}, debug *bool) (*flaggio.Evaluation, error) {
	return nil, errors.ErrNotImplemented
}

func (r *queryResolver) EvaluateAll(ctx context.Context, userID string, context map[string]interface{}, debug *bool) ([]*flaggio.Evaluation, error) {
	return nil, errors.ErrNotImplemented
}
