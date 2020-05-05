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

var _ repository.Evaluation = (*EvaluationRepository)(nil)

// EvaluationRepository implements repository.Evaluation interface using redis.
type EvaluationRepository struct {
	redis *redis.Client
	store repository.Evaluation
	ttl   time.Duration
}

// FindAllByUserID returns all previous flag evaluations for a given user ID.
func (r *EvaluationRepository) FindAllByUserID(ctx context.Context, userID string, search *string, offset, limit *int64) (*flaggio.EvaluationResults, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.FindAllByUserID")
	defer span.Finish()

	// this is called by the admin api, don't cache it
	return r.store.FindAllByUserID(ctx, userID, search, offset, limit)
}

// FindByReqHashAndFlagKey returns a previous flag evaluation for a given request hash and flag key.
func (r *EvaluationRepository) FindByReqHashAndFlagKey(ctx context.Context, reqHash, flagKey string) (*flaggio.Evaluation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.FindByReqHashAndFlagKey")
	defer span.Finish()

	cacheKey := flaggio.EvalCacheKey(reqHash)

	// fetch evaluation results from cache
	cached, err := r.redis.WithContext(ctx).HGet(cacheKey, flagKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		// an unexpected error occurred, return it
		return nil, err
	}
	if cached != "" {
		// cache hit, unmarshal and return result
		var e flaggio.Evaluation
		if err := msgpack.Unmarshal([]byte(cached), &e); err == nil {
			// return if no errors, otherwise defer to the store
			return &e, nil
		}
	}

	// cache miss or disabled, fetch from store
	eval, err := r.store.FindByReqHashAndFlagKey(ctx, reqHash, flagKey)
	if err != nil {
		return nil, err
	}

	// marshal and save result
	if err := r.cacheEvaluations(ctx, reqHash, eval); err != nil {
		return nil, err
	}

	return eval, nil
}

// FindAllByReqHash returns all previous flag evaluations for a given request hash.
func (r *EvaluationRepository) FindAllByReqHash(ctx context.Context, reqHash string) (flaggio.EvaluationList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.FindAllByReqHash")
	defer span.Finish()

	cacheKey := flaggio.EvalCacheKey(reqHash)

	// fetch evaluation results from cache
	cached, err := r.redis.WithContext(ctx).HGetAll(cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		// an unexpected error occurred, return it
		return nil, err
	}
	if len(cached) > 0 {
		// cache hit, unmarshal and return result
		var el flaggio.EvaluationList
		for _, val := range cached {
			var e flaggio.Evaluation
			if err := msgpack.Unmarshal([]byte(val), &e); err == nil {
				el = append(el, &e)
			}
		}
		if len(el) == len(cached) {
			// return if no errors, otherwise defer to the store
			return el, nil
		}
	}

	// cache miss or disabled, fetch from store
	evals, err := r.store.FindAllByReqHash(ctx, reqHash)
	if err != nil {
		return nil, err
	}

	// marshal and save result
	if err := r.cacheEvaluations(ctx, reqHash, evals...); err != nil {
		return nil, err
	}

	return evals, nil
}

// FindByID returns a previous flag evaluation by its ID.
func (r *EvaluationRepository) FindByID(ctx context.Context, id string) (*flaggio.Evaluation, error) {
	// not implemented on cache level
	return r.store.FindByID(ctx, id)
}

// ReplaceOne creates or replaces one evaluation for a user ID.
func (r *EvaluationRepository) ReplaceOne(ctx context.Context, userID string, eval *flaggio.Evaluation) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.ReplaceOne")
	defer span.Finish()

	// call underlying store
	err := r.store.ReplaceOne(ctx, userID, eval)
	if err != nil {
		return err
	}

	// marshal and save result
	return r.cacheEvaluations(ctx, eval.RequestHash, eval)
}

// ReplaceAll creates or replaces evaluations for a combination of user and request hash.
func (r *EvaluationRepository) ReplaceAll(ctx context.Context, userID, reqHash string, evals flaggio.EvaluationList) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.ReplaceAll")
	defer span.Finish()

	// call underlying store
	err := r.store.ReplaceAll(ctx, userID, reqHash, evals)
	if err != nil {
		return err
	}

	// marshal and save result
	return r.cacheEvaluations(ctx, reqHash, evals...)
}

// DeleteAllByUserID deletes evaluations for a user.
func (r *EvaluationRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.DeleteAllByUserID")
	defer span.Finish()

	limit := int64(1)
	res, err := r.store.FindAllByUserID(ctx, userID, nil, nil, &limit)
	if err != nil {
		return err
	}

	if err := r.store.DeleteAllByUserID(ctx, userID); err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx, res.Evaluations[0].RequestHash)
}

// DeleteByID deletes an evaluation by its ID.
func (r *EvaluationRepository) DeleteByID(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.DeleteByID")
	defer span.Finish()

	// find the evaluation
	eval, err := r.store.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// delete the evaluation
	if err := r.store.DeleteByID(ctx, id); err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx, eval.RequestHash, eval.FlagKey)
}

func (r *EvaluationRepository) cacheEvaluations(ctx context.Context, reqHash string, evals ...*flaggio.Evaluation) error {
	if len(evals) == 0 {
		return nil
	}

	redisCtx := r.redis.WithContext(ctx)
	cacheKey := flaggio.EvalCacheKey(reqHash)
	evalsMap := make(map[string]interface{}, len(evals))

	for _, eval := range evals {
		b, err := msgpack.Marshal(eval)
		if err != nil {
			return err
		}
		evalsMap[eval.FlagKey] = b
	}
	if err := redisCtx.HSet(cacheKey, evalsMap).Err(); err != nil {
		return err
	}
	return redisCtx.Expire(cacheKey, r.ttl).Err()
}

func (r *EvaluationRepository) invalidateRelevantCacheKeys(ctx context.Context, reqHash string, flagKeys ...string) error {
	redisCtx := r.redis.WithContext(ctx)
	cacheKey := flaggio.EvalCacheKey(reqHash)

	if len(flagKeys) > 0 {
		// delete all evaluations inside a redis hash
		return redisCtx.HDel(cacheKey, flagKeys...).Err()
	}
	// delete all evaluations by its key
	return redisCtx.Del(cacheKey).Err()
}

// NewEvaluationRepository returns a new evaluation repository that uses redis
// as underlying storage.
func NewEvaluationRepository(redisClient *redis.Client, store repository.Evaluation) repository.Evaluation {
	return &EvaluationRepository{
		redis: redisClient,
		store: store,
		ttl:   1 * time.Hour,
	}
}
