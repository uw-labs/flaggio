package repository

//go:generate mockgen -destination=./mocks/evaluation_mock.go -package=repository_mock github.com/victorkt/flaggio/internal/repository Evaluation

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

// Flag represents a set of operations available to list and manage evaluations.
type Evaluation interface {
	// FindAllByUserID returns all previous flag evaluations for a given user ID.
	FindAllByUserID(ctx context.Context, userID string, search *string, offset, limit *int64) (*flaggio.EvaluationResults, error)
	// FindAllByReqHash returns all previous flag evaluations for a given request hash.
	FindAllByReqHash(ctx context.Context, reqHash string) (flaggio.EvaluationList, error)
	// FindByReqHashAndFlagKey returns a previous flag evaluation for a given request hash and flag key.
	FindByReqHashAndFlagKey(ctx context.Context, reqHash, flagKey string) (*flaggio.Evaluation, error)
	// FindByID returns a previous flag evaluation by its ID.
	FindByID(ctx context.Context, id string) (*flaggio.Evaluation, error)
	// ReplaceOne creates or replaces one evaluation for a user ID.
	ReplaceOne(ctx context.Context, userID string, eval *flaggio.Evaluation) error
	// ReplaceAll creates or replaces evaluations for a combination of user and request hash.
	ReplaceAll(ctx context.Context, userID, reqHash string, evals flaggio.EvaluationList) error
	// DeleteAllByUserID deletes evaluations for a user.
	DeleteAllByUserID(ctx context.Context, userID string) error
	// DeleteByID deletes an evaluation by its ID.
	DeleteByID(ctx context.Context, id string) error
}
