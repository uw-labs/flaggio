package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
)

func TestSegmentRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// drop database first
	if err := mongoDB.Drop(ctx); err != nil {
		t.Fatalf("failed drop database: %s", err)
	}

	// create new repo
	repo, err := mongo_repo.NewSegmentRepository(ctx, mongoDB)
	assert.NoError(t, err, "failed to create segment repository")

	var sgmnt1ID, sgmnt2ID string
	var sgmnt1, sgmnt2 *flaggio.Segment

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "create the first segment",
			run: func(t *testing.T) {
				sgmnt1ID, err = repo.Create(ctx, flaggio.NewSegment{Name: "test users"})
				assert.NoError(t, err, "failed to create first segment")
			},
		},
		{
			name: "checks the segment was created",
			run: func(t *testing.T) {
				sgmnt1, err = repo.FindByID(ctx, sgmnt1ID)
				assert.NoError(t, err, "failed to find first segment")
				assert.Equal(t, newSegment(sgmnt1ID, "test users", sgmnt1.CreatedAt), sgmnt1)
			},
		},
		{
			name: "create the second segment",
			run: func(t *testing.T) {
				sgmnt2ID, err = repo.Create(ctx, flaggio.NewSegment{Name: "beta users"})
				assert.NoError(t, err, "failed to create second segment")
			},
		},
		{
			name: "find the created segment",
			run: func(t *testing.T) {
				sgmnt2, err = repo.FindByID(ctx, sgmnt2ID)
				assert.NoError(t, err, "failed to find second segment")
				assert.Equal(t, newSegment(sgmnt2ID, "beta users", sgmnt2.CreatedAt), sgmnt2)
			},
		},
		{
			name: "find all segments",
			run: func(t *testing.T) {
				sgmnts, err := repo.FindAll(ctx, nil, nil)
				assert.NoError(t, err, "failed to find all segments")
				expectedSegments := []*flaggio.Segment{sgmnt2, sgmnt1} // sorted by name",
				assert.Equal(t, expectedSegments, sgmnts)
			},
		},
		{
			name: "limit segment results",
			run: func(t *testing.T) {
				sgmnts, err := repo.FindAll(ctx, nil, int64Ptr(1))
				assert.NoError(t, err, "failed to limit segment results")
				expectedSegments := []*flaggio.Segment{sgmnt2}
				assert.Equal(t, expectedSegments, sgmnts)
			},
		},
		{
			name: "limit segment results with offset",
			run: func(t *testing.T) {
				sgmnts, err := repo.FindAll(ctx, int64Ptr(1), int64Ptr(1))
				assert.NoError(t, err, "failed to limit segment results with offset")
				expectedSegments := []*flaggio.Segment{sgmnt1}
				assert.Equal(t, expectedSegments, sgmnts)
			},
		},
		{
			name: "update the second segment",
			run: func(t *testing.T) {
				err = repo.Update(ctx, sgmnt2ID, flaggio.UpdateSegment{Name: stringPtr("alpha testers")})
				assert.NoError(t, err, "failed to update second segment")
			},
		},
		{
			name: "checks first segment is untouched",
			run: func(t *testing.T) {
				sgmnt1, err = repo.FindByID(ctx, sgmnt1ID)
				assert.NoError(t, err, "failed to find first segment again")
				assert.Equal(t, newSegment(sgmnt1ID, "test users", sgmnt1.CreatedAt), sgmnt1)
			},
		},
		{
			name: "check second segment was updated",
			run: func(t *testing.T) {
				sgmnt2, err = repo.FindByID(ctx, sgmnt2ID)
				assert.NoError(t, err, "failed to find second segment again")
				expectedSegment := newSegment(sgmnt2ID, "alpha testers", sgmnt2.CreatedAt)
				expectedSegment.UpdatedAt = sgmnt2.UpdatedAt
				assert.Equal(t, expectedSegment, sgmnt2)
				assert.NotNil(t, sgmnt2.UpdatedAt)
			},
		},
		{
			name: "delete the first segment",
			run: func(t *testing.T) {
				err = repo.Delete(ctx, sgmnt1ID)
				assert.NoError(t, err, "failed to delete first segment")
			},
		},
		{
			name: "find deleted segment",
			run: func(t *testing.T) {
				sgmnt1, err = repo.FindByID(ctx, sgmnt1ID)
				assert.EqualError(t, err, "segment: not found")
				assert.Nil(t, sgmnt1)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, tt.run)
	}
}

func newSegment(id, name string, createdAt time.Time) *flaggio.Segment {
	return &flaggio.Segment{
		ID:        id,
		Name:      name,
		Rules:     []*flaggio.SegmentRule{},
		CreatedAt: createdAt,
		UpdatedAt: nil,
	}
}
