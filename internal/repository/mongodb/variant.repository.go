package mongodb

import (
	"context"
	"time"

	"github.com/victorkt/flaggio/internal/errors"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ repository.Variant = VariantRepository{}

// VariantRepository implements repository.Variant interface using mongodb.
type VariantRepository struct {
	flagRepo *FlagRepository
}

// FindByID returns a variant that has a given ID.
func (r VariantRepository) FindByID(ctx context.Context, flagIDHex, idHex string) (*flaggio.Variant, error) {
	flagID, err := primitive.ObjectIDFromHex(flagIDHex)
	if err != nil {
		return nil, err
	}
	variantID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": flagID, "variants._id": variantID}
	projection := bson.M{"variants.$": 1}
	opts := options.FindOne().SetProjection(projection)

	var f flagModel
	if err := r.flagRepo.col.FindOne(ctx, filter, opts).Decode(&f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("variant")
		}
		return nil, err
	}
	if len(f.Variants) != 1 {
		return nil, errors.NotFound("variant")
	}
	return f.Variants[0].asVariant(), nil
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
