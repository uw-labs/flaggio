package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	repository_mock "github.com/victorkt/flaggio/internal/repository/mocks"
	redis_repo "github.com/victorkt/flaggio/internal/repository/redis"
)

var (
	segmentResults = []*flaggio.Segment{
		{ID: "1", Name: "s1"},
		{ID: "2", Name: "s2"},
	}
)

func TestSegmentRepository_FindAll(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying repository on cache miss",
			run: func(t *testing.T, segmmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmmentStoreRepo)
				segmmentStoreRepo.EXPECT().FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, nil).
					Times(1).Return(segmentResults, nil)

				res, err := segmentRedisRepo.FindAll(ctx, nil, nil)
				assert.NoError(t, err)
				assert.Equal(t, segmentResults, res)
			},
		},
		{
			name: "doesnt call underlying repository on cache hit",
			run: func(t *testing.T, segmmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmmentStoreRepo)
				segmmentStoreRepo.EXPECT().FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, nil).
					Times(0)

				res, err := segmentRedisRepo.FindAll(ctx, nil, nil)
				assert.NoError(t, err)
				assert.Equal(t, segmentResults, res)
			},
		},
		{
			name: "always calls underlying repository on offset query",
			run: func(t *testing.T, segmmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmmentStoreRepo)
				segmmentStoreRepo.EXPECT().FindAll(gomock.AssignableToTypeOf(ctxInterface), int64Ptr(1), nil).
					Times(1).Return(segmentResults, nil)

				res, err := segmentRedisRepo.FindAll(ctx, int64Ptr(1), nil)
				assert.NoError(t, err)
				assert.Equal(t, segmentResults, res)
			},
		},
		{
			name: "always calls underlying repository on limit query",
			run: func(t *testing.T, segmmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmmentStoreRepo)
				segmmentStoreRepo.EXPECT().FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, int64Ptr(10)).
					Times(1).Return(segmentResults, nil)

				res, err := segmentRedisRepo.FindAll(ctx, nil, int64Ptr(10))
				assert.NoError(t, err)
				assert.Equal(t, segmentResults, res)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			segmmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)

			tt.run(t, segmmentStoreRepo)
		})
	}
}

func TestSegmentRepository_FindByID(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying repository on cache miss",
			run: func(t *testing.T, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				sgmnt := segmentResults[0]
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmentStoreRepo)
				segmentStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(1).Return(sgmnt, nil)

				res, err := segmentRedisRepo.FindByID(ctx, "1")
				assert.NoError(t, err)
				assert.Equal(t, sgmnt, res)
			},
		},
		{
			name: "doesnt call underlying repository on cache hit",
			run: func(t *testing.T, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				sgmnt := segmentResults[0]
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmentStoreRepo)
				segmentStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(0)

				res, err := segmentRedisRepo.FindByID(ctx, "1")
				assert.NoError(t, err)
				assert.Equal(t, sgmnt, res)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)

			tt.run(t, segmentStoreRepo)
		})
	}
}

func TestSegmentRepository_Create(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "clears all cached evaluations",
			run: func(t *testing.T, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				redisCtx := redisClient.WithContext(ctx)

				// cache an evaluation
				err := redisCtx.Set(flaggio.EvalCacheKey("test"), "whatever", 10*time.Minute).Err()
				assert.NoError(t, err)

				// verify there is a cached evaluation key
				cachedKeys, err := redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)

				// prepare repository mock
				sgmnt := segmentResults[0]
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmentStoreRepo)
				segmentStoreRepo.EXPECT().Create(gomock.AssignableToTypeOf(ctxInterface), flaggio.NewSegment{Name: "s1"}).
					Times(1).Return(sgmnt.ID, nil)

				// call redis repository
				id, err := segmentRedisRepo.Create(ctx, flaggio.NewSegment{Name: "s1"})
				assert.NoError(t, err)
				assert.Equal(t, sgmnt.ID, id)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 0)
			},
		},
		{
			name: "clears relevant cached segments",
			run: func(t *testing.T, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				redisCtx := redisClient.WithContext(ctx)

				// cache a segment
				err := redisCtx.Set(flaggio.SegmentCacheKey("1"), "whatever", 10*time.Minute).Err()
				assert.NoError(t, err)

				// verify there is a cached segment key
				cachedKeys, err := redisCtx.Keys(flaggio.SegmentCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)

				// prepare repository mock
				sgmnt := segmentResults[0]
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmentStoreRepo)
				segmentStoreRepo.EXPECT().Create(gomock.AssignableToTypeOf(ctxInterface), flaggio.NewSegment{Name: "s1"}).
					Times(1).Return(sgmnt.ID, nil)

				// call redis repository
				id, err := segmentRedisRepo.Create(ctx, flaggio.NewSegment{Name: "s1"})
				assert.NoError(t, err)
				assert.Equal(t, sgmnt.ID, id)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.SegmentCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 0)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)

			tt.run(t, segmentStoreRepo)
		})
	}
}

func TestSegmentRepository_Update(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "clears all cached evaluations",
			run: func(t *testing.T, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				redisCtx := redisClient.WithContext(ctx)

				// cache an evaluation
				err := redisCtx.Set(flaggio.EvalCacheKey("test"), "whatever", 10*time.Minute).Err()
				assert.NoError(t, err)

				// verify there is a cached evaluation key
				cachedKeys, err := redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)

				// prepare repository mock
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmentStoreRepo)
				segmentStoreRepo.EXPECT().Update(gomock.AssignableToTypeOf(ctxInterface), "1", flaggio.UpdateSegment{Name: stringPtr("s1")}).
					Times(1).Return(nil)

				// call redis repository
				err = segmentRedisRepo.Update(ctx, "1", flaggio.UpdateSegment{Name: stringPtr("s1")})
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 0)
			},
		},
		{
			name: "clears relevant cached segments",
			run: func(t *testing.T, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				redisCtx := redisClient.WithContext(ctx)

				// cache a segment
				err := redisCtx.Set(flaggio.SegmentCacheKey("1"), "whatever", 10*time.Minute).Err()
				assert.NoError(t, err)

				// verify there is a cached segment key
				cachedKeys, err := redisCtx.Keys(flaggio.SegmentCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)

				// prepare repository mock
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmentStoreRepo)
				segmentStoreRepo.EXPECT().Update(gomock.AssignableToTypeOf(ctxInterface), "1", flaggio.UpdateSegment{Name: stringPtr("s1")}).
					Times(1).Return(nil)

				// call redis repository
				err = segmentRedisRepo.Update(ctx, "1", flaggio.UpdateSegment{Name: stringPtr("s1")})
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.SegmentCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 0)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)

			tt.run(t, segmentStoreRepo)
		})
	}
}

func TestSegmentRepository_Delete(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "clears all cached evaluations",
			run: func(t *testing.T, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				redisCtx := redisClient.WithContext(ctx)

				// cache an evaluation
				err := redisCtx.Set(flaggio.EvalCacheKey("test"), "whatever", 10*time.Minute).Err()
				assert.NoError(t, err)

				// verify there is a cached evaluation key
				cachedKeys, err := redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)

				// prepare repository mock
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmentStoreRepo)
				segmentStoreRepo.EXPECT().Delete(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(1).Return(nil)

				// call redis repository
				err = segmentRedisRepo.Delete(ctx, "1")
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 0)
			},
		},
		{
			name: "clears relevant cached segments",
			run: func(t *testing.T, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				redisCtx := redisClient.WithContext(ctx)

				// cache a segment
				err := redisCtx.Set(flaggio.SegmentCacheKey("1"), "whatever", 10*time.Minute).Err()
				assert.NoError(t, err)

				// verify there is a cached segment key
				cachedKeys, err := redisCtx.Keys(flaggio.SegmentCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)

				// prepare repository mock
				segmentRedisRepo := redis_repo.NewSegmentRepository(redisClient, segmentStoreRepo)
				segmentStoreRepo.EXPECT().Delete(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(1).Return(nil)

				// call redis repository
				err = segmentRedisRepo.Delete(ctx, "1")
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.SegmentCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 0)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)

			tt.run(t, segmentStoreRepo)
		})
	}
}
