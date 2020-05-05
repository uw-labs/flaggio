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
	frl = &flaggio.FlagRule{
		Rule: flaggio.Rule{ID: "1"},
	}
	srl = &flaggio.SegmentRule{
		Rule: flaggio.Rule{ID: "1"},
	}
)

func TestRuleRepository_FindFlagRuleByID(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockRule, *repository_mock.MockFlag, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying repository",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().FindFlagRuleByID(gomock.AssignableToTypeOf(ctxInterface), "2", "1").
					Times(2).Return(frl, nil)

				res, err := ruleRedisRepo.FindFlagRuleByID(ctx, "2", "1")
				assert.NoError(t, err)
				assert.Equal(t, frl, res)

				res2, err2 := ruleRedisRepo.FindFlagRuleByID(ctx, "2", "1")
				assert.NoError(t, err2)
				assert.Equal(t, frl, res2)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)
			ruleStoreRepo := repository_mock.NewMockRule(mockCtrl)

			tt.run(t, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
		})
	}
}

func TestRuleRepository_CreateFlagRule(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockRule, *repository_mock.MockFlag, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "doesn't clear cached evaluations",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
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
				flg := flagResults.Flags[0]
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().CreateFlagRule(gomock.AssignableToTypeOf(ctxInterface), "2", flaggio.NewFlagRule{}).
					Times(1).Return(frl.ID, nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				id, err := ruleRedisRepo.CreateFlagRule(ctx, "2", flaggio.NewFlagRule{})
				assert.NoError(t, err)
				assert.Equal(t, frl.ID, id)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				redisCtx := redisClient.WithContext(ctx)

				// cache a flag
				err := redisCtx.Set(flaggio.FlagCacheKey("key", "f1"), "whatever", 10*time.Minute).Err()
				assert.NoError(t, err)

				// verify there is a cached flag key
				cachedKeys, err := redisCtx.Keys(flaggio.FlagCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)

				// prepare repository mock
				flg := flagResults.Flags[0]
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().CreateFlagRule(gomock.AssignableToTypeOf(ctxInterface), "2", flaggio.NewFlagRule{}).
					Times(1).Return(frl.ID, nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				id, err := ruleRedisRepo.CreateFlagRule(ctx, "2", flaggio.NewFlagRule{})
				assert.NoError(t, err)
				assert.Equal(t, frl.ID, id)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.FlagCacheKey("*")).Result()
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
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)
			ruleStoreRepo := repository_mock.NewMockRule(mockCtrl)

			tt.run(t, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
		})
	}
}

func TestRuleRepository_UpdateFlagRule(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockRule, *repository_mock.MockFlag, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "doesn't clear cached evaluations",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
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
				flg := flagResults.Flags[0]
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().UpdateFlagRule(gomock.AssignableToTypeOf(ctxInterface), "2", "1", flaggio.UpdateFlagRule{}).
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				err = ruleRedisRepo.UpdateFlagRule(ctx, "2", "1", flaggio.UpdateFlagRule{})
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				redisCtx := redisClient.WithContext(ctx)

				// cache a flag
				err := redisCtx.Set(flaggio.FlagCacheKey("key", "f1"), "whatever", 10*time.Minute).Err()
				assert.NoError(t, err)

				// verify there is a cached flag key
				cachedKeys, err := redisCtx.Keys(flaggio.FlagCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)

				// prepare repository mock
				flg := flagResults.Flags[0]
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().UpdateFlagRule(gomock.AssignableToTypeOf(ctxInterface), "2", "1", flaggio.UpdateFlagRule{}).
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				err = ruleRedisRepo.UpdateFlagRule(ctx, "2", "1", flaggio.UpdateFlagRule{})
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.FlagCacheKey("*")).Result()
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
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)
			ruleStoreRepo := repository_mock.NewMockRule(mockCtrl)

			tt.run(t, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
		})
	}
}

func TestRuleRepository_DeleteFlagRule(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockRule, *repository_mock.MockFlag, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "doesn't clear cached evaluations",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
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
				flg := flagResults.Flags[0]
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().DeleteFlagRule(gomock.AssignableToTypeOf(ctxInterface), "2", "1").
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				err = ruleRedisRepo.DeleteFlagRule(ctx, "2", "1")
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				redisCtx := redisClient.WithContext(ctx)

				// cache a flag
				err := redisCtx.Set(flaggio.FlagCacheKey("key", "f1"), "whatever", 10*time.Minute).Err()
				assert.NoError(t, err)

				// verify there is a cached flag key
				cachedKeys, err := redisCtx.Keys(flaggio.FlagCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)

				// prepare repository mock
				flg := flagResults.Flags[0]
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().DeleteFlagRule(gomock.AssignableToTypeOf(ctxInterface), "2", "1").
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				err = ruleRedisRepo.DeleteFlagRule(ctx, "2", "1")
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.FlagCacheKey("*")).Result()
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
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)
			ruleStoreRepo := repository_mock.NewMockRule(mockCtrl)

			tt.run(t, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
		})
	}
}

func TestRuleRepository_FindSegmentRuleByID(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockRule, *repository_mock.MockFlag, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying repository",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().FindSegmentRuleByID(gomock.AssignableToTypeOf(ctxInterface), "2", "1").
					Times(2).Return(srl, nil)

				res, err := ruleRedisRepo.FindSegmentRuleByID(ctx, "2", "1")
				assert.NoError(t, err)
				assert.Equal(t, srl, res)

				res2, err2 := ruleRedisRepo.FindSegmentRuleByID(ctx, "2", "1")
				assert.NoError(t, err2)
				assert.Equal(t, srl, res2)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)
			ruleStoreRepo := repository_mock.NewMockRule(mockCtrl)

			tt.run(t, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
		})
	}
}

func TestRuleRepository_CreateSegmentRule(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockRule, *repository_mock.MockFlag, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "clears all cached evaluations",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
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
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().CreateSegmentRule(gomock.AssignableToTypeOf(ctxInterface), "1", flaggio.NewSegmentRule{}).
					Times(1).Return(srl.ID, nil)

				// call redis repository
				id, err := ruleRedisRepo.CreateSegmentRule(ctx, "1", flaggio.NewSegmentRule{})
				assert.NoError(t, err)
				assert.Equal(t, srl.ID, id)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 0)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
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
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().CreateSegmentRule(gomock.AssignableToTypeOf(ctxInterface), "1", flaggio.NewSegmentRule{}).
					Times(1).Return(srl.ID, nil)

				// call redis repository
				id, err := ruleRedisRepo.CreateSegmentRule(ctx, "1", flaggio.NewSegmentRule{})
				assert.NoError(t, err)
				assert.Equal(t, srl.ID, id)

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
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)
			ruleStoreRepo := repository_mock.NewMockRule(mockCtrl)

			tt.run(t, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
		})
	}
}

func TestRuleRepository_UpdateSegmentRule(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockRule, *repository_mock.MockFlag, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "clears all cached evaluations",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
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
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().UpdateSegmentRule(gomock.AssignableToTypeOf(ctxInterface), "1", "1", flaggio.UpdateSegmentRule{}).
					Times(1).Return(nil)

				// call redis repository
				err = ruleRedisRepo.UpdateSegmentRule(ctx, "1", "1", flaggio.UpdateSegmentRule{})
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 0)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
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
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().UpdateSegmentRule(gomock.AssignableToTypeOf(ctxInterface), "1", "1", flaggio.UpdateSegmentRule{}).
					Times(1).Return(nil)

				// call redis repository
				err = ruleRedisRepo.UpdateSegmentRule(ctx, "1", "1", flaggio.UpdateSegmentRule{})
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
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)
			ruleStoreRepo := repository_mock.NewMockRule(mockCtrl)

			tt.run(t, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
		})
	}
}

func TestRuleRepository_DeleteSegmentRule(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockRule, *repository_mock.MockFlag, *repository_mock.MockSegment)
	}{
		// these tests are meant to be run in order
		{
			name: "clears all cached evaluations",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
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
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().DeleteSegmentRule(gomock.AssignableToTypeOf(ctxInterface), "1", "1").
					Times(1).Return(nil)

				// call redis repository
				err = ruleRedisRepo.DeleteSegmentRule(ctx, "1", "1")
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 0)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, ruleStoreRepo *repository_mock.MockRule, flagStoreRepo *repository_mock.MockFlag, segmentStoreRepo *repository_mock.MockSegment) {
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
				ruleRedisRepo := redis_repo.NewRuleRepository(redisClient, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
				ruleStoreRepo.EXPECT().DeleteSegmentRule(gomock.AssignableToTypeOf(ctxInterface), "1", "1").
					Times(1).Return(nil)

				// call redis repository
				err = ruleRedisRepo.DeleteSegmentRule(ctx, "1", "1")
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
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentStoreRepo := repository_mock.NewMockSegment(mockCtrl)
			ruleStoreRepo := repository_mock.NewMockRule(mockCtrl)

			tt.run(t, ruleStoreRepo, flagStoreRepo, segmentStoreRepo)
		})
	}
}
