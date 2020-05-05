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
	vrnt = &flaggio.Variant{
		ID: "1", Value: 1,
	}
)

func TestVariantRepository_FindByID(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockVariant, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying repository",
			run: func(t *testing.T, variantStoreRepo *repository_mock.MockVariant, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				variantRedisRepo := redis_repo.NewVariantRepository(redisClient, variantStoreRepo, flagStoreRepo)
				variantStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2", "1").
					Times(2).Return(vrnt, nil)

				res, err := variantRedisRepo.FindByID(ctx, "2", "1")
				assert.NoError(t, err)
				assert.Equal(t, vrnt, res)

				res2, err2 := variantRedisRepo.FindByID(ctx, "2", "1")
				assert.NoError(t, err2)
				assert.Equal(t, vrnt, res2)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)
			variantStoreRepo := repository_mock.NewMockVariant(mockCtrl)

			tt.run(t, variantStoreRepo, flagStoreRepo)
		})
	}
}

func TestVariantRepository_Create(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockVariant, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "doesn't clear cached evaluations",
			run: func(t *testing.T, variantStoreRepo *repository_mock.MockVariant, flagStoreRepo *repository_mock.MockFlag) {
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
				variantRedisRepo := redis_repo.NewVariantRepository(redisClient, variantStoreRepo, flagStoreRepo)
				variantStoreRepo.EXPECT().Create(gomock.AssignableToTypeOf(ctxInterface), "2", flaggio.NewVariant{Value: 1}).
					Times(1).Return(vrnt.ID, nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				id, err := variantRedisRepo.Create(ctx, "2", flaggio.NewVariant{Value: 1})
				assert.NoError(t, err)
				assert.Equal(t, vrnt.ID, id)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, variantStoreRepo *repository_mock.MockVariant, flagStoreRepo *repository_mock.MockFlag) {
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
				variantRedisRepo := redis_repo.NewVariantRepository(redisClient, variantStoreRepo, flagStoreRepo)
				variantStoreRepo.EXPECT().Create(gomock.AssignableToTypeOf(ctxInterface), "2", flaggio.NewVariant{Value: false}).
					Times(1).Return(vrnt.ID, nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				id, err := variantRedisRepo.Create(ctx, "2", flaggio.NewVariant{Value: false})
				assert.NoError(t, err)
				assert.Equal(t, vrnt.ID, id)

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
			variantStoreRepo := repository_mock.NewMockVariant(mockCtrl)

			tt.run(t, variantStoreRepo, flagStoreRepo)
		})
	}
}

func TestVariantRepository_Update(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockVariant, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "doesn't clear cached evaluations",
			run: func(t *testing.T, variantStoreRepo *repository_mock.MockVariant, flagStoreRepo *repository_mock.MockFlag) {
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
				variantRedisRepo := redis_repo.NewVariantRepository(redisClient, variantStoreRepo, flagStoreRepo)
				variantStoreRepo.EXPECT().Update(gomock.AssignableToTypeOf(ctxInterface), "2", "1", flaggio.UpdateVariant{Value: "abc"}).
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				err = variantRedisRepo.Update(ctx, "2", "1", flaggio.UpdateVariant{Value: "abc"})
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, variantStoreRepo *repository_mock.MockVariant, flagStoreRepo *repository_mock.MockFlag) {
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
				variantRedisRepo := redis_repo.NewVariantRepository(redisClient, variantStoreRepo, flagStoreRepo)
				variantStoreRepo.EXPECT().Update(gomock.AssignableToTypeOf(ctxInterface), "2", "1", flaggio.UpdateVariant{Value: "cde"}).
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				err = variantRedisRepo.Update(ctx, "2", "1", flaggio.UpdateVariant{Value: "cde"})
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
			variantStoreRepo := repository_mock.NewMockVariant(mockCtrl)

			tt.run(t, variantStoreRepo, flagStoreRepo)
		})
	}
}

func TestVariantRepository_Delete(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockVariant, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "doesn't clear cached evaluations",
			run: func(t *testing.T, variantStoreRepo *repository_mock.MockVariant, flagStoreRepo *repository_mock.MockFlag) {
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
				variantRedisRepo := redis_repo.NewVariantRepository(redisClient, variantStoreRepo, flagStoreRepo)
				variantStoreRepo.EXPECT().Delete(gomock.AssignableToTypeOf(ctxInterface), "2", "1").
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				err = variantRedisRepo.Delete(ctx, "2", "1")
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, variantStoreRepo *repository_mock.MockVariant, flagStoreRepo *repository_mock.MockFlag) {
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
				variantRedisRepo := redis_repo.NewVariantRepository(redisClient, variantStoreRepo, flagStoreRepo)
				variantStoreRepo.EXPECT().Delete(gomock.AssignableToTypeOf(ctxInterface), "2", "1").
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "2").
					Times(1).Return(flg, nil)

				// call redis repository
				err = variantRedisRepo.Delete(ctx, "2", "1")
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
			variantStoreRepo := repository_mock.NewMockVariant(mockCtrl)

			tt.run(t, variantStoreRepo, flagStoreRepo)
		})
	}
}
