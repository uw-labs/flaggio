package api

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

var _ QueryResolver = &queryResolver{}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Ping(ctx context.Context) (*bool, error) {
	pong := true
	return &pong, nil
}

func (r *queryResolver) Evaluate(
	ctx context.Context,
	flagKey, userID string,
	usrContext map[string]interface{},
	debug *bool,
) (*flaggio.Evaluation, error) {
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

	res, err := flaggio.Evaluate(usrContext, flg)
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

func (r *queryResolver) EvaluateAll(
	ctx context.Context,
	userID string,
	usrContext map[string]interface{},
) ([]*flaggio.Evaluation, error) {
	flgs, err := r.flagRepo.FindAll(ctx, nil, nil)
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

	// TODO: use go routines
	var evals []*flaggio.Evaluation
	for _, flg := range flgs {
		flg.Populate(iders)

		res, err := flaggio.Evaluate(usrContext, flg)
		if err != nil {
			return nil, err
		}

		evals = append(evals, &flaggio.Evaluation{
			FlagKey: flg.Key,
			Value:   res.Answer,
			Version: flg.Version,
		})
	}
	return evals, nil
}
