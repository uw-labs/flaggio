package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
)

func TestFlagRuleRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// drop database first
	if err := mongoDB.Drop(ctx); err != nil {
		t.Fatalf("failed drop database: %s", err)
	}

	// create new repo
	flgRepo, err := mongo_repo.NewFlagRepository(ctx, mongoDB)
	assert.NoError(t, err, "failed to create flag repository")
	vrntRepo := mongo_repo.NewVariantRepository(flgRepo.(*mongo_repo.FlagRepository))
	repo := mongo_repo.NewRuleRepository(flgRepo.(*mongo_repo.FlagRepository), nil)

	// create a flag and variant
	flgID, err := flgRepo.Create(ctx, flaggio.NewFlag{Key: "test"})
	assert.NoError(t, err, "failed to create flag")
	vrntID, err := vrntRepo.Create(ctx, flgID, flaggio.NewVariant{Value: "abc"})
	assert.NoError(t, err, "failed to create variant")

	var rl1ID string

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "create a rule",
			run: func(t *testing.T) {
				rl1ID, err = repo.CreateFlagRule(ctx, flgID, flaggio.NewFlagRule{
					Constraints:   []*flaggio.NewConstraint{{Operation: flaggio.OperationOneOf, Property: "name", Values: []interface{}{"test"}}},
					Distributions: []*flaggio.NewDistribution{{VariantID: vrntID, Percentage: 100}},
				})
				assert.NoError(t, err, "failed to create rule")
			},
		},
		{
			name: "checks the rule was created",
			run: func(t *testing.T) {
				rl, err := repo.FindFlagRuleByID(ctx, flgID, rl1ID)
				assert.NoError(t, err, "failed to find rule")
				assert.Equal(t, &flaggio.FlagRule{
					Rule: flaggio.Rule{
						ID: rl.ID, // use the generated id
						Constraints: []*flaggio.Constraint{
							{ID: rl.Constraints[0].ID, Operation: flaggio.OperationOneOf, Property: "name", Values: []interface{}{"test"}},
						},
					},
					Distributions: []*flaggio.Distribution{
						{ID: rl.Distributions[0].ID, Variant: &flaggio.Variant{ID: vrntID, Value: "abc"}, Percentage: 100},
					},
				}, rl)
			},
		},
		{
			name: "update the rule",
			run: func(t *testing.T) {
				err := repo.UpdateFlagRule(ctx, flgID, rl1ID, flaggio.UpdateFlagRule{
					Constraints:   []*flaggio.NewConstraint{{Operation: flaggio.OperationGreater, Property: "age", Values: []interface{}{18}}},
					Distributions: []*flaggio.NewDistribution{{VariantID: vrntID, Percentage: 50}},
				})
				assert.NoError(t, err, "failed to update rule")
			},
		},
		{
			name: "check updated rule",
			run: func(t *testing.T) {
				rl, err := repo.FindFlagRuleByID(ctx, flgID, rl1ID)
				assert.NoError(t, err, "failed to find updated rule")
				assert.Equal(t, &flaggio.FlagRule{
					Rule: flaggio.Rule{
						ID: rl.ID, // use the generated id
						Constraints: []*flaggio.Constraint{
							{ID: rl.Constraints[0].ID, Property: "age", Operation: flaggio.OperationGreater, Values: []interface{}{int32(18)}},
						},
					},
					Distributions: []*flaggio.Distribution{
						{ID: rl.Distributions[0].ID, Variant: &flaggio.Variant{ID: vrntID, Value: "abc"}, Percentage: 50},
					},
				}, rl)
			},
		},
		{
			name: "delete the rule",
			run: func(t *testing.T) {
				err := repo.DeleteFlagRule(ctx, flgID, rl1ID)
				assert.NoError(t, err, "failed to delete rule")
			},
		},
		{
			name: "find deleted rule",
			run: func(t *testing.T) {
				rl, err := repo.FindFlagRuleByID(ctx, flgID, rl1ID)
				assert.EqualError(t, err, "rule: not found")
				assert.Nil(t, rl)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, tt.run)
	}
}

func TestSegmentRuleRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// drop database first
	if err := mongoDB.Drop(ctx); err != nil {
		t.Fatalf("failed drop database: %s", err)
	}

	// create new repo
	sgmntRepo, err := mongo_repo.NewSegmentRepository(ctx, mongoDB)
	assert.NoError(t, err, "failed to create flag repository")
	repo := mongo_repo.NewRuleRepository(nil, sgmntRepo.(*mongo_repo.SegmentRepository))

	// create a segment
	sgmntID, err := sgmntRepo.Create(ctx, flaggio.NewSegment{Name: "beta testers"})
	assert.NoError(t, err, "failed to create segment")

	var rl1ID string

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "create a rule",
			run: func(t *testing.T) {
				rl1ID, err = repo.CreateSegmentRule(ctx, sgmntID, flaggio.NewSegmentRule{
					Constraints: []*flaggio.NewConstraint{{Operation: flaggio.OperationOneOf, Property: "name", Values: []interface{}{"test"}}},
				})
				assert.NoError(t, err, "failed to create rule")
			},
		},
		{
			name: "checks the rule was created",
			run: func(t *testing.T) {
				rl, err := repo.FindSegmentRuleByID(ctx, sgmntID, rl1ID)
				assert.NoError(t, err, "failed to find rule")
				assert.Equal(t, &flaggio.SegmentRule{
					Rule: flaggio.Rule{
						ID: rl.ID, // use the generated id
						Constraints: []*flaggio.Constraint{
							{ID: rl.Constraints[0].ID, Operation: flaggio.OperationOneOf, Property: "name", Values: []interface{}{"test"}},
						},
					},
				}, rl)
			},
		},
		{
			name: "update the rule",
			run: func(t *testing.T) {
				err := repo.UpdateSegmentRule(ctx, sgmntID, rl1ID, flaggio.UpdateSegmentRule{
					Constraints: []*flaggio.NewConstraint{{Operation: flaggio.OperationGreater, Property: "age", Values: []interface{}{18}}},
				})
				assert.NoError(t, err, "failed to update rule")
			},
		},
		{
			name: "check updated rule",
			run: func(t *testing.T) {
				rl, err := repo.FindSegmentRuleByID(ctx, sgmntID, rl1ID)
				assert.NoError(t, err, "failed to find updated rule")
				assert.Equal(t, &flaggio.SegmentRule{
					Rule: flaggio.Rule{
						ID: rl.ID, // use the generated id
						Constraints: []*flaggio.Constraint{
							{ID: rl.Constraints[0].ID, Operation: flaggio.OperationGreater, Property: "age", Values: []interface{}{int32(18)}},
						},
					},
				}, rl)
			},
		},
		{
			name: "delete the rule",
			run: func(t *testing.T) {
				err := repo.DeleteSegmentRule(ctx, sgmntID, rl1ID)
				assert.NoError(t, err, "failed to delete rule")
			},
		},
		{
			name: "find deleted rule",
			run: func(t *testing.T) {
				rl, err := repo.FindSegmentRuleByID(ctx, sgmntID, rl1ID)
				assert.EqualError(t, err, "rule: not found")
				assert.Nil(t, rl)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, tt.run)
	}
}
