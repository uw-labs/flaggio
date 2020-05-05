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

var _ repository.Segment = (*SegmentRepository)(nil)

// SegmentRepository implements repository.Segment interface using redis.
type SegmentRepository struct {
	redis *redis.Client
	store repository.Segment
	ttl   time.Duration
}

// FindAll returns a list of segments, based on an optional offset and limit.
func (r *SegmentRepository) FindAll(ctx context.Context, offset, limit *int64) ([]*flaggio.Segment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisSegmentRepository.FindAll")
	defer span.Finish()

	shouldCache := shouldCacheFindAll(nil, offset, limit)
	cacheKey := flaggio.SegmentCacheKey("*")

	if shouldCache {
		// fetch results from cache
		cached, err := r.redis.WithContext(ctx).Get(cacheKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			// an unexpected error occurred, return it
			return nil, err
		}
		if cached != "" {
			// cache hit, unmarshal and return result
			var s []*flaggio.Segment
			if err := msgpack.Unmarshal([]byte(cached), &s); err == nil {
				// return if no errors, otherwise defer to the store
				return s, nil
			}
		}
	}

	// cache miss or disabled, fetch from store
	res, err := r.store.FindAll(ctx, offset, limit)
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

// FindByID returns a segment that has a given ID.
func (r *SegmentRepository) FindByID(ctx context.Context, id string) (*flaggio.Segment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisSegmentRepository.FindByID")
	defer span.Finish()

	cacheKey := flaggio.SegmentCacheKey(id)

	// fetch results from cache
	cached, err := r.redis.WithContext(ctx).Get(cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		// an unexpected error occurred, return it
		return nil, err
	}
	if cached != "" {
		// cache hit, unmarshal and return result
		var s flaggio.Segment
		if err := msgpack.Unmarshal([]byte(cached), &s); err == nil {
			// return if no errors, otherwise defer to the store
			return &s, nil
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

// Create creates a new segment.
func (r *SegmentRepository) Create(ctx context.Context, input flaggio.NewSegment) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisSegmentRepository.Create")
	defer span.Finish()

	id, err := r.store.Create(ctx, input)
	if err != nil {
		return "", err
	}

	// invalidate all relevant keys
	return id, r.invalidateRelevantCacheKeys(ctx, id)
}

// Update updates a segment.
func (r *SegmentRepository) Update(ctx context.Context, id string, input flaggio.UpdateSegment) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisSegmentRepository.Update")
	defer span.Finish()

	if err := r.store.Update(ctx, id, input); err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx, id)
}

// Delete deletes a segment.
func (r *SegmentRepository) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisSegmentRepository.Delete")
	defer span.Finish()

	if err := r.store.Delete(ctx, id); err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx, id)
}

func (r *SegmentRepository) invalidateRelevantCacheKeys(ctx context.Context, segmentID string) error {
	redisCtx := r.redis.WithContext(ctx)

	// invalidate all relevant keys
	keysToInvalidate, err := redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
	if err != nil {
		return err
	}
	keysToInvalidate = append(
		[]string{
			flaggio.SegmentCacheKey("*"),
			flaggio.SegmentCacheKey(segmentID),
		},
		keysToInvalidate...,
	)

	return redisCtx.Del(keysToInvalidate...).Err()
}

// NewSegmentRepository returns a new segment repository that uses redis
// as underlying storage.
func NewSegmentRepository(redisClient *redis.Client, store repository.Segment) repository.Segment {
	return &SegmentRepository{
		redis: redisClient,
		store: store,
		ttl:   1 * time.Hour,
	}
}
