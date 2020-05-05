package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
)

func TestVariantRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// drop database first
	if err := mongoDB.Drop(ctx); err != nil {
		t.Fatalf("failed drop database: %s", err)
	}

	// create new repo
	flgRepo, err := mongo_repo.NewFlagRepository(ctx, mongoDB)
	assert.NoError(t, err, "failed to create flag repository")
	repo := mongo_repo.NewVariantRepository(flgRepo.(*mongo_repo.FlagRepository))

	// create a flag
	flgID, err := flgRepo.Create(ctx, flaggio.NewFlag{Key: "test"})
	assert.NoError(t, err, "failed to create flag")

	var vrnt1ID, vrnt2ID string

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "create the first variant",
			run: func(t *testing.T) {
				vrnt1ID, err = repo.Create(ctx, flgID, flaggio.NewVariant{Value: 2.1})
				assert.NoError(t, err, "failed to create first variant")
			},
		},
		{
			name: "checks the variant was created",
			run: func(t *testing.T) {
				vrnt, err := repo.FindByID(ctx, flgID, vrnt1ID)
				assert.NoError(t, err, "failed to find first variant")
				assert.Equal(t, &flaggio.Variant{ID: vrnt1ID, Value: 2.1}, vrnt)
			},
		},
		{
			name: "create the second variant",
			run: func(t *testing.T) {
				vrnt2ID, err = repo.Create(ctx, flgID, flaggio.NewVariant{Value: "a"})
				assert.NoError(t, err, "failed to create second variant")
			},
		},
		{
			name: "find the created variant",
			run: func(t *testing.T) {
				vrnt, err := repo.FindByID(ctx, flgID, vrnt2ID)
				assert.NoError(t, err, "failed to find second variant")
				assert.Equal(t, &flaggio.Variant{ID: vrnt2ID, Value: "a"}, vrnt)
			},
		},
		{
			name: "update the second variant",
			run: func(t *testing.T) {
				err := repo.Update(ctx, flgID, vrnt2ID, flaggio.UpdateVariant{Value: false})
				assert.NoError(t, err, "failed to update second variant")
			},
		},
		{
			name: "find second variant",
			run: func(t *testing.T) {
				vrnt, err := repo.FindByID(ctx, flgID, vrnt2ID)
				assert.NoError(t, err, "failed to find second variant again")
				assert.Equal(t, &flaggio.Variant{ID: vrnt2ID, Value: false}, vrnt)
			},
		},
		{
			name: "delete the first variant",
			run: func(t *testing.T) {
				err := repo.Delete(ctx, flgID, vrnt1ID)
				assert.NoError(t, err, "failed to delete first variant")
			},
		},
		{
			name: "find deleted variant",
			run: func(t *testing.T) {
				vrnt, err := repo.FindByID(ctx, flgID, vrnt1ID)
				assert.EqualError(t, err, "variant: not found")
				assert.Nil(t, vrnt)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, tt.run)
	}

}
