package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
)

func TestEvaluationRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// drop database first
	if err := mongoDB.Drop(ctx); err != nil {
		t.Fatalf("failed drop database: %s", err)
	}

	// create new repos
	flgRepo, err := mongo_repo.NewFlagRepository(ctx, mongoDB)
	assert.NoError(t, err, "failed to create flag repository")
	repo, err := mongo_repo.NewEvaluationRepository(ctx, mongoDB)
	assert.NoError(t, err, "failed to create evaluation repository")

	// create two flags
	flg1ID, err := flgRepo.Create(ctx, flaggio.NewFlag{Key: "test"})
	assert.NoError(t, err, "failed to create first flag")
	flg1, err := flgRepo.FindByID(ctx, flg1ID)
	assert.NoError(t, err, "failed to find first flag")
	flg2ID, err := flgRepo.Create(ctx, flaggio.NewFlag{Key: "height"})
	assert.NoError(t, err, "failed to create second flag")
	flg2, err := flgRepo.FindByID(ctx, flg2ID)
	assert.NoError(t, err, "failed to find second flag")

	// create evaluations
	evaluations := []*flaggio.Evaluation{
		{FlagID: flg2.ID, FlagKey: flg2.Key, FlagVersion: flg2.Version, RequestHash: "123456789", Value: "abc"},
		{FlagID: flg1.ID, FlagKey: flg1.Key, FlagVersion: flg1.Version, RequestHash: "123456789", Value: 2.1},
	}

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		// these tests are meant to be run in order
		{
			name: "create evaluations for user",
			run: func(t *testing.T) {
				err := repo.ReplaceAll(ctx, "TEST1", "123456789", evaluations)
				assert.NoError(t, err, "failed to create first evaluations")
			},
		},
		{
			name: "evaluations were created correctly",
			run: func(t *testing.T) {
				evals, err := repo.FindAllByUserID(ctx, "TEST1", nil, nil, nil)
				assert.NoError(t, err, "failed to find first evaluations")
				evaluations[0].ID = evals.Evaluations[0].ID
				evaluations[0].CreatedAt = evals.Evaluations[0].CreatedAt
				evaluations[1].ID = evals.Evaluations[1].ID
				evaluations[1].CreatedAt = evals.Evaluations[1].CreatedAt
				assert.Equal(t, &flaggio.EvaluationResults{Evaluations: evaluations, Total: 2}, evals)
			},
		},
		{
			name: "search for evaluation",
			run: func(t *testing.T) {
				evals, err := repo.FindAllByUserID(ctx, "TEST1", stringPtr("height"), nil, nil)
				assert.NoError(t, err, "failed to find search evaluations")
				assert.Equal(t, &flaggio.EvaluationResults{Evaluations: []*flaggio.Evaluation{evaluations[0]}, Total: 1}, evals)
			},
		},
		{
			name: "find all by user with limit",
			run: func(t *testing.T) {
				evals, err := repo.FindAllByUserID(ctx, "TEST1", nil, nil, int64Ptr(1))
				assert.NoError(t, err, "failed to find all evaluations with limit")
				assert.Equal(t, &flaggio.EvaluationResults{Evaluations: []*flaggio.Evaluation{evaluations[0]}, Total: 2}, evals)
			},
		},
		{
			name: "find all by user with limit and offset",
			run: func(t *testing.T) {
				evals, err := repo.FindAllByUserID(ctx, "TEST1", nil, int64Ptr(1), int64Ptr(1))
				assert.NoError(t, err, "failed to find all evaluations with limit and offset")
				assert.Equal(t, &flaggio.EvaluationResults{Evaluations: []*flaggio.Evaluation{evaluations[1]}, Total: 2}, evals)
			},
		},
		{
			name: "find all by user and flag key",
			run: func(t *testing.T) {
				eval, err := repo.FindByReqHashAndFlagKey(ctx, "123456789", flg1.Key)
				assert.NoError(t, err, "failed to find all by hash and flag id")
				assert.Equal(t, evaluations[1], eval)
			},
		},
		{
			name: "update evaluations for user",
			run: func(t *testing.T) {
				evaluations[0].RequestHash = "aaaaaabbbbbb"
				evaluations[0].Value = false
				err = repo.ReplaceAll(ctx, "TEST1", "123456789", evaluations)
				assert.NoError(t, err, "failed to update evaluations for user")
			},
		},
		{
			name: "checks the evaluations were updated",
			run: func(t *testing.T) {
				evals, err := repo.FindAllByUserID(ctx, "TEST1", nil, nil, nil)
				assert.NoError(t, err, "failed to find updated evaluations")
				evaluations[0].ID = evals.Evaluations[0].ID
				evaluations[0].CreatedAt = evals.Evaluations[0].CreatedAt
				evaluations[1].ID = evals.Evaluations[1].ID
				evaluations[1].CreatedAt = evals.Evaluations[1].CreatedAt
				assert.Equal(t, &flaggio.EvaluationResults{Evaluations: evaluations, Total: 2}, evals)
			},
		},
		{
			name: "delete evaluation by id",
			run: func(t *testing.T) {
				err = repo.DeleteByID(ctx, evaluations[1].ID)
				assert.NoError(t, err, "failed to delete evaluation by id")
			},
		},
		{
			name: "check evaluation was deleted",
			run: func(t *testing.T) {
				eval, err := repo.FindByReqHashAndFlagKey(ctx, "123456789", flg1.Key)
				assert.EqualError(t, err, "evaluation: not found")
				assert.Nil(t, eval)
			},
		},
		{
			name: "delete evaluation by id",
			run: func(t *testing.T) {
				err = repo.DeleteAllByUserID(ctx, "TEST1")
				assert.NoError(t, err, "failed to delete all evaluations by user id")
			},
		},
		{
			name: "check all evaluations were deleted",
			run: func(t *testing.T) {
				evals, err := repo.FindAllByUserID(ctx, "TEST1", nil, nil, nil)
				assert.NoError(t, err, "failed to check all evaluations were deleted")
				assert.Equal(t, &flaggio.EvaluationResults{Total: 0}, evals)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, tt.run)
	}
}
