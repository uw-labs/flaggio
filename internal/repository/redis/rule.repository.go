package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
)

var _ repository.Rule = (*RuleRepository)(nil)

// RuleRepository implements repository.Rule interface using redis.
type RuleRepository struct {
	redis        *redis.Client
	store        repository.Rule
	flagStore    repository.Flag
	segmentStore repository.Segment
	ttl          time.Duration
}

// FindFlagRuleByID returns a flag rule that has a given ID.
func (r *RuleRepository) FindFlagRuleByID(ctx context.Context, flagIDHex, idHex string) (*flaggio.FlagRule, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRuleRepository.FindFlagRuleByID")
	defer span.Finish()

	// no caching for rules
	return r.store.FindFlagRuleByID(ctx, flagIDHex, idHex)
}

// CreateFlagRule creates a new rule under a flag.
func (r *RuleRepository) CreateFlagRule(ctx context.Context, flagID string, input flaggio.NewFlagRule) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRuleRepository.CreateFlagRule")
	defer span.Finish()

	id, err := r.store.CreateFlagRule(ctx, flagID, input)
	if err != nil {
		return "", err
	}

	// invalidate all relevant keys
	return id, r.invalidateFlagRelevantCacheKeys(ctx, flagID)
}

// UpdateFlagRule updates a rule under a flag.
func (r *RuleRepository) UpdateFlagRule(ctx context.Context, flagID, id string, input flaggio.UpdateFlagRule) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRuleRepository.UpdateFlagRule")
	defer span.Finish()

	err := r.store.UpdateFlagRule(ctx, flagID, id, input)
	if err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateFlagRelevantCacheKeys(ctx, flagID)
}

// DeleteFlagRule deletes a rule under a flag.
func (r *RuleRepository) DeleteFlagRule(ctx context.Context, flagID, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRuleRepository.DeleteFlagRule")
	defer span.Finish()

	err := r.store.DeleteFlagRule(ctx, flagID, id)
	if err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateFlagRelevantCacheKeys(ctx, flagID)
}

// FindSegmentRuleByID returns a segment rule that has a given ID.
func (r *RuleRepository) FindSegmentRuleByID(ctx context.Context, segmentIDHex, idHex string) (*flaggio.SegmentRule, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRuleRepository.FindSegmentRuleByID")
	defer span.Finish()

	// no caching for rules
	return r.store.FindSegmentRuleByID(ctx, segmentIDHex, idHex)
}

// CreateSegmentRule creates a new rule under a segment.
func (r *RuleRepository) CreateSegmentRule(ctx context.Context, segmentID string, input flaggio.NewSegmentRule) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRuleRepository.CreateSegmentRule")
	defer span.Finish()

	id, err := r.store.CreateSegmentRule(ctx, segmentID, input)
	if err != nil {
		return "", err
	}

	// invalidate all relevant keys
	return id, r.invalidateSegmentRelevantCacheKeys(ctx, segmentID)
}

// UpdateSegmentRule updates a rule under a segment.
func (r *RuleRepository) UpdateSegmentRule(ctx context.Context, segmentID, id string, input flaggio.UpdateSegmentRule) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRuleRepository.UpdateSegmentRule")
	defer span.Finish()

	err := r.store.UpdateSegmentRule(ctx, segmentID, id, input)
	if err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateSegmentRelevantCacheKeys(ctx, segmentID)
}

// DeleteSegmentRule deletes a rule under a segment.
func (r *RuleRepository) DeleteSegmentRule(ctx context.Context, segmentID, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRuleRepository.DeleteSegmentRule")
	defer span.Finish()

	err := r.store.DeleteSegmentRule(ctx, segmentID, id)
	if err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateSegmentRelevantCacheKeys(ctx, segmentID)
}

func (r *RuleRepository) invalidateFlagRelevantCacheKeys(ctx context.Context, flagID string) error {
	// find the flag so we can get the flag key
	f, err := r.flagStore.FindByID(ctx, flagID)
	if err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.redis.WithContext(ctx).Del(
		flaggio.FlagCacheKey("*"),
		flaggio.FlagCacheKey(flagID),
		flaggio.FlagCacheKey("key", f.Key),
	).Err()
}

func (r *RuleRepository) invalidateSegmentRelevantCacheKeys(ctx context.Context, segmentID string) error {
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

// NewRuleRepository returns a new rule repository that uses redis
// as underlying storage.
func NewRuleRepository(redisClient *redis.Client, store repository.Rule, flagStore repository.Flag, segmentStore repository.Segment) repository.Rule {
	return &RuleRepository{
		redis:        redisClient,
		store:        store,
		flagStore:    flagStore,
		segmentStore: segmentStore,
		ttl:          1 * time.Hour,
	}
}
