package service

//go:generate mockgen -destination=./mocks/flag_mock.go -package=service_mock github.com/uw-labs/flaggio/internal/service Flag

import (
	"context"
)

// Flag holds the logic for evaluating flags
type Flag interface {
	// Evaluate returns the result of an evaluation of a single flag.
	Evaluate(ctx context.Context, flagKey string, req *EvaluationRequest) (*EvaluationResponse, error)
	// EvaluateAll returns the results of the evaluation of all flags.
	EvaluateAll(ctx context.Context, req *EvaluationRequest) (*EvaluationsResponse, error)
}
