package service

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
)

var _ Flag = (*flagService)(nil)

// NewFlagService returns a new Flag
func NewFlagService(
	flagsRepo repository.Flag,
	segmentsRepo repository.Segment,
) Flag {
	return &flagService{
		flagsRepo:    flagsRepo,
		segmentsRepo: segmentsRepo,
	}
}

type flagService struct {
	flagsRepo    repository.Flag
	segmentsRepo repository.Segment
}

// Evaluate evaluates a flag by key, returning a value based on the user context
func (s *flagService) Evaluate(ctx context.Context, flagKey string, req *EvaluationRequest) (*EvaluationResponse, error) {
	flg, err := s.flagsRepo.FindByKey(ctx, flagKey)
	if err != nil {
		return nil, err
	}
	sgmts, err := segmentsAsIdentifiers(s.segmentsRepo.FindAll(ctx, nil, nil))
	if err != nil {
		return nil, err
	}

	flg.Populate(sgmts)

	res, err := flaggio.Evaluate(req.UserContext, flg)
	if err != nil {
		return nil, err
	}

	evalRes := &EvaluationResponse{
		Evaluation: &flaggio.Evaluation{
			FlagKey: flagKey,
			Value:   res.Answer,
		},
	}

	if req.Debug != nil && *req.Debug {
		evalRes.Evaluation.StackTrace = res.Stack()
		evalRes.UserContext = &req.UserContext
	}

	return evalRes, nil
}

// EvaluateAll evaluates all flags, returning a value or an error for each flag based on the user context
func (s *flagService) EvaluateAll(ctx context.Context, req *EvaluationRequest) (*EvaluationsResponse, error) {
	flgs, err := s.flagsRepo.FindAll(ctx, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	sgmts, err := segmentsAsIdentifiers(s.segmentsRepo.FindAll(ctx, nil, nil))
	if err != nil {
		return nil, err
	}

	evals := make([]*flaggio.Evaluation, flgs.Total)
	for idx, flg := range flgs.Flags {
		flg.Populate(sgmts)

		evltn := &flaggio.Evaluation{
			FlagKey: flg.Key,
		}
		res, err := flaggio.Evaluate(req.UserContext, flg)
		if err != nil {
			evltn.Error = err.Error()
		} else {
			evltn.Value = res.Answer
		}

		evals[idx] = evltn
	}

	evalRes := &EvaluationsResponse{
		Evaluations: evals,
	}

	if req.Debug != nil && *req.Debug {
		evalRes.UserContext = &req.UserContext
	}

	return evalRes, nil
}

func segmentsAsIdentifiers(sgmts []*flaggio.Segment, err error) ([]flaggio.Identifier, error) {
	if err != nil {
		return nil, err
	}
	iders := make([]flaggio.Identifier, len(sgmts))
	for idx, sgmnt := range sgmts {
		iders[idx] = sgmnt
	}
	return iders, nil
}
