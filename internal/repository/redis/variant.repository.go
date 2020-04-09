package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
)

var _ repository.Variant = (*VariantRepository)(nil)

// VariantRepository implements repository.Variant interface using redis.
type VariantRepository struct {
	redis     *redis.Client
	store     repository.Variant
	flagStore repository.Flag
	ttl       time.Duration
}

// FindByID returns a variant that has a given ID.
func (r *VariantRepository) FindByID(ctx context.Context, flagIDHex, idHex string) (*flaggio.Variant, error) {
	// no caching for variants
	return r.store.FindByID(ctx, flagIDHex, idHex)
}

// Create creates a new variant.
func (r *VariantRepository) Create(ctx context.Context, flagID string, input flaggio.NewVariant) (string, error) {
	id, err := r.store.Create(ctx, flagID, input)
	if err != nil {
		return "", err
	}

	// invalidate all relevant keys
	return id, r.invalidateRelevantCacheKeys(ctx, flagID)
}

// Update updates a variant.
func (r *VariantRepository) Update(ctx context.Context, flagID, id string, input flaggio.UpdateVariant) error {
	if err := r.store.Update(ctx, flagID, id, input); err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx, flagID)
}

// Delete deletes a variant.
func (r *VariantRepository) Delete(ctx context.Context, flagID, id string) error {
	// delete the flag
	if err := r.store.Delete(ctx, flagID, id); err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx, flagID)
}

func (r *VariantRepository) invalidateRelevantCacheKeys(ctx context.Context, flagID string) error {
	// find the flag so we can get the flag key
	f, err := r.flagStore.FindByID(ctx, flagID)
	if err != nil {
		return err
	}

	redisCtx := r.redis.WithContext(ctx)

	// invalidate all relevant keys
	keysToInvalidate, err := redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
	if err != nil {
		return err
	}
	keysToInvalidate = append(
		[]string{
			flaggio.FlagCacheKey("*"),
			flaggio.FlagCacheKey(flagID),
			flaggio.FlagCacheKey("key", f.Key),
		},
		keysToInvalidate...,
	)

	return redisCtx.Del(keysToInvalidate...).Err()
}

// NewVariantRepository returns a new variant repository that uses redis
// as underlying storage.
func NewVariantRepository(redisClient *redis.Client, store repository.Variant, flagStore repository.Flag) repository.Variant {
	return &VariantRepository{
		redis:     redisClient,
		store:     store,
		flagStore: flagStore,
		ttl:       1 * time.Hour,
	}
}
