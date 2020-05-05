package redis_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	repository_mock "github.com/victorkt/flaggio/internal/repository/mocks"
	redis_repo "github.com/victorkt/flaggio/internal/repository/redis"
)

var (
	flagResults = &flaggio.FlagResults{
		Flags: []*flaggio.Flag{
			{ID: "1", Key: "f1", Enabled: true},
			{ID: "2", Key: "f2", Enabled: false},
		},
		Total: 2,
	}
	ctxInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func TestFlagRepository_FindAll(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying repository on cache miss",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, nil, nil).
					Times(1).Return(flagResults, nil)

				res, err := flagRedisRepo.FindAll(ctx, nil, nil, nil)
				assert.NoError(t, err)
				assert.Equal(t, flagResults, res)
			},
		},
		{
			name: "doesnt call underlying repository on cache hit",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, nil, nil).
					Times(0)

				res, err := flagRedisRepo.FindAll(ctx, nil, nil, nil)
				assert.NoError(t, err)
				assert.Equal(t, flagResults, res)
			},
		},
		{
			name: "always calls underlying repository when searching",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().FindAll(gomock.AssignableToTypeOf(ctxInterface), stringPtr("my search"), nil, nil).
					Times(1).Return(flagResults, nil)

				res, err := flagRedisRepo.FindAll(ctx, stringPtr("my search"), nil, nil)
				assert.NoError(t, err)
				assert.Equal(t, flagResults, res)
			},
		},
		{
			name: "always calls underlying repository on offset query",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, int64Ptr(1), nil).
					Times(1).Return(flagResults, nil)

				res, err := flagRedisRepo.FindAll(ctx, nil, int64Ptr(1), nil)
				assert.NoError(t, err)
				assert.Equal(t, flagResults, res)
			},
		},
		{
			name: "always calls underlying repository on limit query",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, nil, int64Ptr(10)).
					Times(1).Return(flagResults, nil)

				res, err := flagRedisRepo.FindAll(ctx, nil, nil, int64Ptr(10))
				assert.NoError(t, err)
				assert.Equal(t, flagResults, res)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)

			tt.run(t, flagStoreRepo)
		})
	}
}

func TestFlagRepository_FindByID(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying repository on cache miss",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				flg := flagResults.Flags[0]
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(1).Return(flg, nil)

				res, err := flagRedisRepo.FindByID(ctx, "1")
				assert.NoError(t, err)
				assert.Equal(t, flg, res)
			},
		},
		{
			name: "doesnt call underlying repository on cache hit",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				flg := flagResults.Flags[0]
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(0)

				res, err := flagRedisRepo.FindByID(ctx, "1")
				assert.NoError(t, err)
				assert.Equal(t, flg, res)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)

			tt.run(t, flagStoreRepo)
		})
	}
}

func TestFlagRepository_FindByKey(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying repository on cache miss",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				flg := flagResults.Flags[1]
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().FindByKey(gomock.AssignableToTypeOf(ctxInterface), "f2").
					Times(1).Return(flg, nil)

				res, err := flagRedisRepo.FindByKey(ctx, "f2")
				assert.NoError(t, err)
				assert.Equal(t, flg, res)
			},
		},
		{
			name: "doesnt call underlying repository on cache hit",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				flg := flagResults.Flags[1]
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().FindByKey(gomock.AssignableToTypeOf(ctxInterface), "f2").
					Times(0)

				res, err := flagRedisRepo.FindByKey(ctx, "f2")
				assert.NoError(t, err)
				assert.Equal(t, flg, res)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			flagStoreRepo := repository_mock.NewMockFlag(mockCtrl)

			tt.run(t, flagStoreRepo)
		})
	}
}

func TestFlagRepository_Create(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "doesn't clear cached evaluations",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
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
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().Create(gomock.AssignableToTypeOf(ctxInterface), flaggio.NewFlag{Key: "f1", Name: "f1"}).
					Times(1).Return(flg.ID, nil)

				// call redis repository
				id, err := flagRedisRepo.Create(ctx, flaggio.NewFlag{Key: "f1", Name: "f1"})
				assert.NoError(t, err)
				assert.Equal(t, flg.ID, id)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
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
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().Create(gomock.AssignableToTypeOf(ctxInterface), flaggio.NewFlag{Key: "f1", Name: "f1"}).
					Times(1).Return(flg.ID, nil)

				// call redis repository
				id, err := flagRedisRepo.Create(ctx, flaggio.NewFlag{Key: "f1", Name: "f1"})
				assert.NoError(t, err)
				assert.Equal(t, flg.ID, id)

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

			tt.run(t, flagStoreRepo)
		})
	}
}

func TestFlagRepository_Update(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "doesn't clear cached evaluations & searches for flag when flag key is not provided",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
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
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().Update(gomock.AssignableToTypeOf(ctxInterface), "1", flaggio.UpdateFlag{Name: stringPtr("f1")}).
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(1).Return(flg, nil)

				// call redis repository
				err = flagRedisRepo.Update(ctx, "1", flaggio.UpdateFlag{Name: stringPtr("f1")})
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)
			},
		},
		{
			name: "clears relevant cached flags & uses provided flag key",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
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
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().Update(gomock.AssignableToTypeOf(ctxInterface), "1", flaggio.UpdateFlag{Key: stringPtr("f1")}).
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(0)

				// call redis repository
				err = flagRedisRepo.Update(ctx, "1", flaggio.UpdateFlag{Key: stringPtr("f1")})
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

			tt.run(t, flagStoreRepo)
		})
	}
}

func TestFlagRepository_Delete(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T, *repository_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "doesn't clear cached evaluations",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
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
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().Delete(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(1).Return(flg, nil)

				// call redis repository
				err = flagRedisRepo.Delete(ctx, "1")
				assert.NoError(t, err)

				// check cached keys are cleared
				cachedKeys, err = redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
				assert.NoError(t, err)
				assert.Len(t, cachedKeys, 1)
			},
		},
		{
			name: "clears relevant cached flags",
			run: func(t *testing.T, flagStoreRepo *repository_mock.MockFlag) {
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
				flagRedisRepo := redis_repo.NewFlagRepository(redisClient, flagStoreRepo)
				flagStoreRepo.EXPECT().Delete(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(1).Return(nil)
				flagStoreRepo.EXPECT().FindByID(gomock.AssignableToTypeOf(ctxInterface), "1").
					Times(1).Return(flg, nil)

				// call redis repository
				err = flagRedisRepo.Delete(ctx, "1")
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

			tt.run(t, flagStoreRepo)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}
