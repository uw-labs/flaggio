package service

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
)

// FlagService holds the logic for evaluating flags
type FlagService interface {
	Evaluate(ctx context.Context, flagKey string, req *EvaluationRequest) (*flaggio.Evaluation, error)
	EvaluateAll(ctx context.Context, req *EvaluationRequest) (flaggio.EvaluationList, error)
}

// NewFlagService returns a new FlagService
func NewFlagService(
	flagsRepo repository.Flag,
	segmentsRepo repository.Segment,
) FlagService {
	return &flagService{
		flagsRepo:    flagsRepo,
		segmentsRepo: segmentsRepo,
	}
}

type flagService struct {
	flagsRepo    repository.Flag
	segmentsRepo repository.Segment
	router       chi.Router
}

// Evaluate evaluates a flag by key, returning a value based on the user context
func (s *flagService) Evaluate(ctx context.Context, flagKey string, req *EvaluationRequest) (*flaggio.Evaluation, error) {
	flg, err := s.flagsRepo.FindByKey(ctx, flagKey)
	if err != nil {
		return nil, err
	}
	sgmts, err := s.segmentsRepo.FindAll(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	iders := make([]flaggio.Identifier, len(sgmts))
	for idx, sgmnt := range sgmts {
		iders[idx] = sgmnt
	}

	flg.Populate(iders)

	res, err := flaggio.Evaluate(req.UserContext, flg)
	if err != nil {
		return nil, err
	}

	eval := &flaggio.Evaluation{
		FlagKey: flagKey,
		Value:   res.Answer,
	}

	if req.Debug != nil && *req.Debug {
		eval.StackTrace = res.Stack()
	}

	return eval, nil
}

// EvaluateAll evaluates all flags, returning a value or an error for each flag based on the user context
func (s *flagService) EvaluateAll(ctx context.Context, req *EvaluationRequest) (flaggio.EvaluationList, error) {
	flgs, err := s.flagsRepo.FindAll(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	sgmts, err := s.segmentsRepo.FindAll(ctx, nil, nil)
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

		evltn := &flaggio.Evaluation{
			FlagKey: flg.Key,
		}
		res, err := flaggio.Evaluate(req.UserContext, flg)
		if err != nil {
			evltn.Error = err.Error()
		} else {
			evltn.Value = res.Answer
		}

		evals = append(evals, evltn)
	}
	return evals, nil
}

// EvaluationRequest is the evaluation request object
type EvaluationRequest struct {
	UserID      string              `json:"userId"`
	UserContext flaggio.UserContext `json:"context"`
	Debug       *bool               `json:"debug,omitempty"`
}

// Bind adds additional data to the EvaluationRequest.
// Some special fields are added to the user context:
// * $userId is the user ID provided in the request
// * $ip is the network address that sent the request. if the
//     x-forwarded-for header is defined, it will override the
//     value received from RemoteAddr in the request object
func (er EvaluationRequest) Bind(r *http.Request) error {
	// enrich user context
	er.UserContext["$userId"] = er.UserID
	er.UserContext["$ip"] = r.RemoteAddr
	if ip := r.Header.Get("x-forwarded-for"); ip != "" {
		er.UserContext["$ip"] = ip
	}
	return nil
}
