package api

import (
	"context"

	"github.com/victorkohl/flaggio/internal/errors"
	"github.com/victorkohl/flaggio/internal/flaggio"
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

func (r *queryResolver) EvaluateAll(ctx context.Context, userID string, context map[string]interface{}, debug *bool) ([]*flaggio.Evaluation, error) {
	return nil, errors.ErrNotImplemented
}
