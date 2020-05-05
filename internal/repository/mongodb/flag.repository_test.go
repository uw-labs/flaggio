package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
)

func TestFlagRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// drop database first
	if err := mongoDB.Drop(ctx); err != nil {
		t.Fatalf("failed drop database: %s", err)
	}

	// create new repo
	repo, err := mongo_repo.NewFlagRepository(ctx, mongoDB)
	assert.NoError(t, err, "failed to create flag repository")

	var flg1ID, flg2ID string
	var flg1, flg2 *flaggio.Flag

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		// these tests are meant to be run in order
		{
			name: "create the first flag",
			run: func(t *testing.T) {
				flg1ID, err = repo.Create(ctx, flaggio.NewFlag{Key: "test", Name: "testing"})
				assert.NoError(t, err, "failed to create first flag")
			},
		},
		{
			name: "checks the flag was created",
			run: func(t *testing.T) {
				flg1, err = repo.FindByID(ctx, flg1ID)
				assert.NoError(t, err, "failed to find first flag")
				assert.Equal(t, newFlag(flg1ID, "test", "testing", flg1.CreatedAt), flg1)
			},
		},
		{
			name: "find flag by key",
			run: func(t *testing.T) {
				flg1, err = repo.FindByKey(ctx, "test")
				assert.NoError(t, err, "failed to find first flag by key")
				assert.Equal(t, newFlag(flg1ID, "test", "testing", flg1.CreatedAt), flg1)
			},
		},
		{
			name: "create the second flag",
			run: func(t *testing.T) {
				flg2ID, err = repo.Create(ctx, flaggio.NewFlag{Key: "height", Name: "component height"})
				assert.NoError(t, err, "failed to create second flag")
			},
		},
		{
			name: "find the created flag",
			run: func(t *testing.T) {
				flg2, err = repo.FindByID(ctx, flg2ID)
				assert.NoError(t, err, "failed to find second flag")
				assert.Equal(t, newFlag(flg2ID, "height", "component height", flg2.CreatedAt), flg2)
			},
		},
		{
			name: "find all flags",
			run: func(t *testing.T) {
				flgs, err := repo.FindAll(ctx, nil, nil, nil)
				assert.NoError(t, err, "failed to find all flags")
				expectedFlags := &flaggio.FlagResults{
					Flags: []*flaggio.Flag{flg2, flg1}, // sorted by key
					Total: 2,
				}
				assert.Equal(t, expectedFlags, flgs)
			},
		},
		{
			name: "search flags",
			run: func(t *testing.T) {
				flgs, err := repo.FindAll(ctx, stringPtr("height"), nil, nil)
				assert.NoError(t, err, "failed to search flags")
				expectedFlags := &flaggio.FlagResults{Flags: []*flaggio.Flag{flg2}, Total: 1}
				assert.Equal(t, expectedFlags, flgs)
			},
		},
		{
			name: "limit flag results",
			run: func(t *testing.T) {
				flgs, err := repo.FindAll(ctx, nil, nil, int64Ptr(1))
				assert.NoError(t, err, "failed to limit flag results")
				expectedFlags := &flaggio.FlagResults{Flags: []*flaggio.Flag{flg2}, Total: 2}
				assert.Equal(t, expectedFlags, flgs)
			},
		},
		{
			name: "limit flag results with offset",
			run: func(t *testing.T) {
				flgs, err := repo.FindAll(ctx, nil, int64Ptr(1), int64Ptr(1))
				assert.NoError(t, err, "failed to limit flag results with offset")
				expectedFlags := &flaggio.FlagResults{Flags: []*flaggio.Flag{flg1}, Total: 2}
				assert.Equal(t, expectedFlags, flgs)
			},
		},
		{
			name: "update the second flag",
			run: func(t *testing.T) {
				err := repo.Update(ctx, flg2ID, flaggio.UpdateFlag{Name: stringPtr("button height"), Enabled: boolPtr(true)})
				assert.NoError(t, err, "failed to update second flag")
			},
		},
		{
			name: "checks first flag is untouched",
			run: func(t *testing.T) {
				flg1, err = repo.FindByID(ctx, flg1ID)
				assert.NoError(t, err, "failed to find first flag again")
				assert.Equal(t, newFlag(flg1ID, "test", "testing", flg1.CreatedAt), flg1)
			},
		},
		{
			name: "check second flag was updated",
			run: func(t *testing.T) {
				flg2, err = repo.FindByID(ctx, flg2ID)
				assert.NoError(t, err, "failed to find second flag again")
				expectedFlag := newFlag(flg2ID, "height", "button height", flg2.CreatedAt)
				expectedFlag.Enabled = true
				expectedFlag.Version = 2
				expectedFlag.UpdatedAt = flg2.UpdatedAt
				assert.Equal(t, expectedFlag, flg2)
				assert.NotNil(t, flg2.UpdatedAt)
			},
		},
		{
			name: "delete the first flag",
			run: func(t *testing.T) {
				err := repo.Delete(ctx, flg1ID)
				assert.NoError(t, err, "failed to delete first flag")
			},
		},
		{
			name: "find deleted flag",
			run: func(t *testing.T) {
				flg1, err = repo.FindByID(ctx, flg1ID)
				assert.EqualError(t, err, "flag: not found")
				assert.Nil(t, flg1)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, tt.run)
	}
}

func newFlag(id, key, name string, createdAt time.Time) *flaggio.Flag {
	return &flaggio.Flag{
		ID:                    id,
		Key:                   key,
		Name:                  name,
		Enabled:               false,
		Version:               1,
		Variants:              []*flaggio.Variant{},
		Rules:                 []*flaggio.FlagRule{},
		DefaultVariantWhenOn:  nil,
		DefaultVariantWhenOff: nil,
		CreatedAt:             createdAt,
		UpdatedAt:             nil,
	}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
