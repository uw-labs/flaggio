package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/service"
	"github.com/vmihailenco/msgpack/v4"
)

var _ service.Flag = (*flagService)(nil)

// flagService implements service.Flag interface using redis.
type flagService struct {
	redis *redis.Client
	svc   service.Flag
	ttl   time.Duration
}

// Evaluate returns the result of an evaluation of a single flag.
func (s flagService) Evaluate(ctx context.Context, flagKey string, req *service.EvaluationRequest) (*service.EvaluationResponse, error) {
	shouldCache := shouldCacheEvaluation(req)
	var cacheKey string

	if shouldCache {
		reqHash, err := req.Hash()
		if err != nil {
			return nil, err
		}
		cacheKey = flaggio.EvalCacheKey(flagKey, reqHash)

		// fetch flag results from cache
		cached, err := s.redis.WithContext(ctx).Get(cacheKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			// an unexpected error occurred, return it
			return nil, err
		}
		if cached != "" {
			// cache hit, unmarshall and return result
			var er service.EvaluationResponse
			if err := msgpack.Unmarshal([]byte(cached), &er); err == nil {
				// return if no errors, otherwise defer to the store
				return &er, nil
			}
		}
	}

	// cache miss or debug request, call underlying service
	res, err := s.svc.Evaluate(ctx, flagKey, req)
	if err != nil {
		return nil, err
	}

	if shouldCache {
		// marshall and save result
		b, err := msgpack.Marshal(res)
		if err != nil {
			return nil, err
		}
		if err := s.redis.Set(cacheKey, b, s.ttl).Err(); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// EvaluateAll returns the results of the evaluation of all flags.
func (s flagService) EvaluateAll(ctx context.Context, req *service.EvaluationRequest) (*service.EvaluationsResponse, error) {
	shouldCache := shouldCacheEvaluation(req)
	var cacheKey string

	if shouldCache {
		reqHash, err := req.Hash()
		if err != nil {
			return nil, err
		}
		cacheKey = flaggio.EvalCacheKey(reqHash)

		// fetch flag results from cache
		cached, err := s.redis.WithContext(ctx).Get(cacheKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			// an unexpected error occurred, return it
			return nil, err
		}
		if cached != "" {
			// cache hit, unmarshall and return result
			var er service.EvaluationsResponse
			if err := msgpack.Unmarshal([]byte(cached), &er); err == nil {
				// return if no errors, otherwise defer to the store
				return &er, nil
			}
		}
	}

	// cache miss or debug request, call underlying service
	res, err := s.svc.EvaluateAll(ctx, req)
	if err != nil {
		return nil, err
	}

	if shouldCache {
		// marshall and save result
		b, err := msgpack.Marshal(res)
		if err != nil {
			return nil, err
		}
		if err := s.redis.Set(cacheKey, b, s.ttl).Err(); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func NewFlagService(redisClient *redis.Client, svc service.Flag) service.Flag {
	return &flagService{
		redis: redisClient,
		svc:   svc,
		ttl:   24 * time.Hour,
	}
}
