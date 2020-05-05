package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
	"github.com/vmihailenco/msgpack/v4"
)

var _ repository.Flag = (*FlagRepository)(nil)

// FlagRepository implements repository.Flag interface using redis.
type FlagRepository struct {
	redis *redis.Client
	store repository.Flag
	ttl   time.Duration
}

// FindAll returns a list of flags, based on an optional offset and limit.
func (r *FlagRepository) FindAll(ctx context.Context, search *string, offset, limit *int64) (*flaggio.FlagResults, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisFlagRepository.FindAll")
	defer span.Finish()

	shouldCache := shouldCacheFindAll(search, offset, limit)
	cacheKey := flaggio.FlagCacheKey("*")

	if shouldCache {
		// fetch flag results from cache
		cached, err := r.redis.WithContext(ctx).Get(cacheKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			// an unexpected error occurred, return it
			return nil, err
		}
		if cached != "" {
			// cache hit, unmarshal and return result
			var fr flaggio.FlagResults
			if err := msgpack.Unmarshal([]byte(cached), &fr); err == nil {
				// return if no errors, otherwise defer to the store
				return &fr, nil
			}
		}
	}

	// cache miss or disabled, fetch from store
	res, err := r.store.FindAll(ctx, search, offset, limit)
	if err != nil {
		return nil, err
	}

	// marshal and save result
	if shouldCache {
		b, err := msgpack.Marshal(res)
		if err != nil {
			return nil, err
		}
		if err := r.redis.Set(cacheKey, b, r.ttl).Err(); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// FindByID returns a flag that has a given ID.
func (r *FlagRepository) FindByID(ctx context.Context, id string) (*flaggio.Flag, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisFlagRepository.FindByID")
	defer span.Finish()

	cacheKey := flaggio.FlagCacheKey(id)

	// fetch flag results from cache
	cached, err := r.redis.WithContext(ctx).Get(cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		// an unexpected error occurred, return it
		return nil, err
	}
	if cached != "" {
		// cache hit, unmarshal and return result
		var f flaggio.Flag
		if err := msgpack.Unmarshal([]byte(cached), &f); err == nil {
			// return if no errors, otherwise defer to the store
			return &f, nil
		}
	}

	// cache miss, fetch from store
	res, err := r.store.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// marshal and save result
	b, err := msgpack.Marshal(res)
	if err != nil {
		return nil, err
	}
	if err := r.redis.WithContext(ctx).Set(cacheKey, b, r.ttl).Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// FindByKey returns a flag that has a given key.
func (r *FlagRepository) FindByKey(ctx context.Context, key string) (*flaggio.Flag, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisFlagRepository.FindByKey")
	defer span.Finish()

	cacheKey := flaggio.FlagCacheKey("key", key)

	// fetch flag results from cache
	cached, err := r.redis.WithContext(ctx).Get(cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		// an unexpected error occurred, return it
		return nil, err
	}
	if cached != "" {
		// cache hit, unmarshal and return result
		var f flaggio.Flag
		if err := msgpack.Unmarshal([]byte(cached), &f); err == nil {
			// return if no errors, otherwise defer to the store
			return &f, nil
		}
	}

	// cache miss, fetch from store
	res, err := r.store.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	// marshal and save result
	b, err := msgpack.Marshal(res)
	if err != nil {
		return nil, err
	}
	if err := r.redis.WithContext(ctx).Set(cacheKey, b, r.ttl).Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Create creates a new flag.
func (r *FlagRepository) Create(ctx context.Context, input flaggio.NewFlag) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisFlagRepository.Create")
	defer span.Finish()

	id, err := r.store.Create(ctx, input)
	if err != nil {
		return "", err
	}

	// invalidate all relevant keys
	return id, r.invalidateRelevantCacheKeys(ctx, id, input.Key)
}

// Update updates a flag.
func (r *FlagRepository) Update(ctx context.Context, id string, input flaggio.UpdateFlag) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisFlagRepository.Update")
	defer span.Finish()

	if err := r.store.Update(ctx, id, input); err != nil {
		return err
	}

	var flagKey string
	if input.Key != nil {
		flagKey = *input.Key
	} else {
		f, err := r.FindByID(ctx, id)
		if err != nil {
			return err
		}
		flagKey = f.Key
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx, id, flagKey)
}

// Delete deletes a flag.
func (r *FlagRepository) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisFlagRepository.Delete")
	defer span.Finish()

	// find the flag so we can get the flag key
	f, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// delete the flag
	if err := r.store.Delete(ctx, id); err != nil {
		return nil
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx, id, f.Key)
}

func (r *FlagRepository) invalidateRelevantCacheKeys(ctx context.Context, flagID, flagKey string) error {
	// invalidate all relevant keys
	return r.redis.WithContext(ctx).Del(
		flaggio.FlagCacheKey("*"),
		flaggio.FlagCacheKey(flagID),
		flaggio.FlagCacheKey("key", flagKey),
	).Err()
}

// NewFlagRepository returns a new flag repository that uses redis
// as underlying storage.
func NewFlagRepository(redisClient *redis.Client, store repository.Flag) repository.Flag {
	return &FlagRepository{
		redis: redisClient,
		store: store,
		ttl:   1 * time.Hour,
	}
}
