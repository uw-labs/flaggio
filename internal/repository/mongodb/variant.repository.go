package mongodb

import (
	"context"
	"time"

	"github.com/victorkt/flaggio/internal/errors"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ repository.Variant = VariantRepository{}

// VariantRepository implements repository.Variant interface using mongodb.
type VariantRepository struct {
	flagRepo *FlagRepository
}

// Create creates a new variant under a flag.
func (r VariantRepository) Create(ctx context.Context, flagIDHex string, v flaggio.NewVariant) (string, error) {
	vrntModel := &variantModel{
		ID:          primitive.NewObjectID(),
		Description: v.Description,
		Value:       v.Value,
	}
	flagID, err := primitive.ObjectIDFromHex(flagIDHex)
	if err != nil {
		return "", err
	}
	filter := bson.M{"_id": flagID}
	res, err := r.flagRepo.col.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{"variants": vrntModel},
		"$set":  bson.M{"updatedAt": time.Now()},
		"$inc":  bson.M{"version": 1},
	})
	if err != nil {
		return "", err
	}
	if res.ModifiedCount == 0 {
		return "", errors.NotFound("flag")
	}
	return vrntModel.ID.Hex(), nil
}

// Update updates a variant under a flag.
func (r VariantRepository) Update(ctx context.Context, flagIDHex, idHex string, v flaggio.UpdateVariant) error {
	flagID, err := primitive.ObjectIDFromHex(flagIDHex)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	mods := bson.M{
		"updatedAt": time.Now(),
	}
	if v.Description != nil {
		mods["variants.$.description"] = *v.Description
	}
	if v.Value != nil {
		mods["variants.$.value"] = v.Value
	}
	res, err := r.flagRepo.col.UpdateOne(
		ctx,
		bson.M{"_id": flagID, "variants._id": id},
		bson.M{"$set": mods, "$inc": bson.M{"version": 1}},
	)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.NotFound("variant")
	}
	return nil
}

// Delete deletes a variant under a flag.
func (r VariantRepository) Delete(ctx context.Context, flagIDHex, idHex string) error {
	flagID, err := primitive.ObjectIDFromHex(flagIDHex)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	res, err := r.flagRepo.col.UpdateOne(ctx, bson.M{"_id": flagID}, bson.M{
		"$pull": bson.M{"variants": bson.M{"_id": id}},
		"$set":  bson.M{"updatedAt": time.Now()},
		"$inc":  bson.M{"version": 1},
	})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.NotFound("variant")
	}
	return nil
}

// NewMongoVariantRepository returns a new variant repository that uses mongodb
// as underlying storage.
func NewMongoVariantRepository(flagRepo *FlagRepository) *VariantRepository {
	return &VariantRepository{
		flagRepo: flagRepo,
	}
}
