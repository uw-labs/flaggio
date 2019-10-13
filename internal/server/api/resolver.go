package api

import (
	"context"

	"github.com/victorkohl/flaggio/internal/errors"
	"github.com/victorkohl/flaggio/internal/flaggio"
	"github.com/victorkohl/flaggio/internal/repository"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	flagRepo    repository.Flag
	segmentRepo repository.Segment
}

func NewResolver(
	flagRepo repository.Flag,
	segmentRepo repository.Segment,
) *Resolver {
	return &Resolver{
		flagRepo:    flagRepo,
		segmentRepo: segmentRepo,
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

func (r *queryResolver) Evaluate(ctx context.Context, flagKey string, userID string, usrCtx map[string]interface{}, debug *bool) (*flaggio.Evaluation, error) {
	flg, err := r.flagRepo.FindByKey(ctx, flagKey)
	if err != nil {
		return nil, err
	}
	sgmts, err := r.segmentRepo.FindAll(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	iders := make([]flaggio.Identifier, len(sgmts))
	for idx, sgmnt := range sgmts {
		iders[idx] = sgmnt
	}

	flg.Populate(iders)

	res, err := flaggio.Evaluate(usrCtx, flg)
	if err != nil {
		return nil, err
	}

	eval := &flaggio.Evaluation{
		FlagKey: flagKey,
		Value:   res.Answer,
		Version: flg.Version,
	}

	if debug != nil && *debug {
		eval.StackTrace = res.Stack()
	}

	return eval, nil
}

func (r *queryResolver) EvaluateAll(ctx context.Context, userID string, context map[string]interface{}, debug *bool) ([]*flaggio.Evaluation, error) {
	return nil, errors.ErrNotImplemented
}
