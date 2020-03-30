package service

import (
	"context"
	"net/http"

	"github.com/victorkt/clientip"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
)

// FlagService holds the logic for evaluating flags
type FlagService interface {
	Evaluate(ctx context.Context, flagKey string, req *EvaluationRequest) (*EvaluationResponse, error)
	EvaluateAll(ctx context.Context, req *EvaluationRequest) (*EvaluationsResponse, error)
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
}

// Evaluate evaluates a flag by key, returning a value based on the user context
func (s *flagService) Evaluate(ctx context.Context, flagKey string, req *EvaluationRequest) (*EvaluationResponse, error) {
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

	evalRes := &EvaluationsResponse{
		Evaluations: evals,
	}

	if req.Debug != nil && *req.Debug {
		evalRes.UserContext = &req.UserContext
	}

	return evalRes, nil
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
// * $ip is the network address that originated the request
func (er EvaluationRequest) Bind(r *http.Request) error {
	// enrich user context
	er.UserContext["$userId"] = er.UserID
	er.UserContext["$ip"] = clientip.FromContext(r.Context()).String()
	return nil
}

// EvaluationResponse is the evaluation response object
type EvaluationResponse struct {
	Evaluation  *flaggio.Evaluation  `json:"evaluation"`
	UserContext *flaggio.UserContext `json:"context,omitempty"`
}

// Render can enrich the EvaluationResponse object before being returned to the
// user. Currently it does nothing, but is needed to satisfy the
// chi.Renderer interface.
func (e *EvaluationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// EvaluationsResponse is the evaluation response object
type EvaluationsResponse struct {
	Evaluations flaggio.EvaluationList `json:"evaluations"`
	UserContext *flaggio.UserContext   `json:"context,omitempty"`
}

// Render can enrich the EvaluationsResponse object before being returned to the
// user. Currently it does nothing, but is needed to satisfy the
// chi.Renderer interface.
func (e *EvaluationsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
