package service

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	apperrors "github.com/victorkt/flaggio/internal/errors"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
)

var _ Flag = (*flagService)(nil)

// NewFlagService returns a new Flag service
func NewFlagService(
	flagsRepo repository.Flag,
	segmentsRepo repository.Segment,
	evalsRepo repository.Evaluation,
	usersRepo repository.User,
) Flag {
	return &flagService{
		flagsRepo:    flagsRepo,
		segmentsRepo: segmentsRepo,
		evalsRepo:    evalsRepo,
		usersRepo:    usersRepo,
	}
}

type flagService struct {
	flagsRepo    repository.Flag
	segmentsRepo repository.Segment
	evalsRepo    repository.Evaluation
	usersRepo    repository.User
}

// Evaluate evaluates a flag by key, returning a value based on the user context
func (s *flagService) Evaluate(ctx context.Context, flagKey string, req *EvaluationRequest) (*EvaluationResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlagService.Evaluate")
	defer span.Finish()

	// fetch flag
	flg, err := s.flagsRepo.FindByKey(ctx, flagKey)
	if err != nil {
		return nil, err
	}
	// fetch previous evaluations for this flag
	hash, err := req.Hash()
	if err != nil {
		return nil, err
	}
	eval, err := s.evalsRepo.FindByReqHashAndFlagKey(ctx, hash, flg.Key)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	// if there are no previous evaluations, evaluate the flag
	if invalidEval(hash, flg, eval) {
		sgmts, err := segmentsAsIdentifiers(s.segmentsRepo.FindAll(ctx, nil, nil))
		if err != nil {
			return nil, err
		}

		flg.Populate(sgmts)

		evalSpan, _ := opentracing.StartSpanFromContext(ctx, "flaggio.Evaluate")
		res, err := flaggio.Evaluate(req.UserContext, flg)
		evalSpan.Finish()
		if err != nil {
			return nil, err
		}

		eval = &flaggio.Evaluation{
			FlagID:      flg.ID,
			FlagVersion: flg.Version,
			FlagKey:     flg.Key,
			RequestHash: hash,
			Value:       res.Answer,
		}
		if req.IsDebug() {
			eval.StackTrace = res.Stack()
		} else {
			// create or update the user
			if err := s.usersRepo.Replace(ctx, req.UserID, req.UserContext); err != nil {
				return nil, err
			}
			// replace the evaluation for the user and flag key
			if err := s.evalsRepo.ReplaceOne(ctx, req.UserID, eval); err != nil {
				return nil, err
			}
		}
	}

	// build the response
	evalRes := &EvaluationResponse{
		Evaluation: eval,
	}

	if req.IsDebug() {
		evalRes.UserContext = &req.UserContext
	}

	return evalRes, nil
}

// EvaluateAll evaluates all flags, returning a value or an error for each flag based on the user context
func (s *flagService) EvaluateAll(ctx context.Context, req *EvaluationRequest) (*EvaluationsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlagService.EvaluateAll")
	defer span.Finish()

	// fetch previous evaluations
	hash, err := req.Hash()
	if err != nil {
		return nil, err
	}
	prevEvals, err := s.evalsRepo.FindAllByReqHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	// fetch all flags
	flgs, err := s.flagsRepo.FindAll(ctx, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	// fetch segments
	sgmts, err := segmentsAsIdentifiers(s.segmentsRepo.FindAll(ctx, nil, nil))
	if err != nil {
		return nil, err
	}

	// check for missing flag evaluations
	validEvals := validFlagEvals(hash, flgs.Flags, prevEvals)
	evals := make(flaggio.EvaluationList, len(flgs.Flags))

	// evaluate flags
	evalSpan, _ := opentracing.StartSpanFromContext(ctx, "flaggio.Evaluate")
	var outdatedEvals flaggio.EvaluationList
	for idx, flg := range flgs.Flags {
		if evltn, ok := validEvals[flg.ID]; ok {
			evals[idx] = evltn
			continue
		}
		flg.Populate(sgmts)

		evltn := &flaggio.Evaluation{
			FlagID:      flg.ID,
			FlagVersion: flg.Version,
			FlagKey:     flg.Key,
			RequestHash: hash,
		}
		res, err := flaggio.Evaluate(req.UserContext, flg)
		if err != nil {
			evltn.Error = err.Error()
		} else {
			evltn.Value = res.Answer
			outdatedEvals = append(outdatedEvals, evltn)
		}

		evals[idx] = evltn
	}
	evalSpan.Finish()

	if len(outdatedEvals) > 0 && !req.IsDebug() {
		// create or update the user
		if err := s.usersRepo.Replace(ctx, req.UserID, req.UserContext); err != nil {
			return nil, err
		}
		// replace the evaluations for the user
		if err := s.evalsRepo.ReplaceAll(ctx, req.UserID, hash, outdatedEvals); err != nil {
			return nil, err
		}
	}

	// build the response
	evalRes := &EvaluationsResponse{
		Evaluations: evals,
	}

	if req.IsDebug() {
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

func invalidEval(reqHash string, flg *flaggio.Flag, eval *flaggio.Evaluation) bool {
	// eval is invalid if no evaluation found, the evaluation was
	// for a previous flag version or the user context changed
	return flg == nil ||
		eval == nil ||
		flg.Version != eval.FlagVersion ||
		reqHash != eval.RequestHash
}

func validFlagEvals(reqHash string, flgs []*flaggio.Flag, evals flaggio.EvaluationList) map[string]*flaggio.Evaluation {
	validEvals := map[string]*flaggio.Evaluation{}
	for _, eval := range evals {
		validEvals[eval.FlagID] = eval
	}
	for _, flg := range flgs {
		eval := validEvals[flg.ID]
		if invalidEval(reqHash, flg, eval) {
			delete(validEvals, flg.ID)
		}
	}
	return validEvals
}
