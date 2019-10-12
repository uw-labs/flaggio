package mongodb

import (
	"context"

	"github.com/victorkohl/flaggio/internal/errors"
	"github.com/victorkohl/flaggio/internal/flaggio"
	"github.com/victorkohl/flaggio/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	_ repository.Variant = &VariantRepository{}
)

type VariantRepository struct {
	flagRepo *FlagRepository
}

func (r VariantRepository) Create(ctx context.Context, flagIDHex string, v flaggio.NewVariant) (*flaggio.Variant, error) {
	var defaultWhenOn, defaultWhenOff bool
	if v.DefaultWhenOn != nil {
		defaultWhenOn = *v.DefaultWhenOn
	}
	if v.DefaultWhenOff != nil {
		defaultWhenOff = *v.DefaultWhenOff
	}
	vrntModel := &variantModel{
		ID:             primitive.NewObjectID(),
		Key:            v.Key,
		Description:    v.Description,
		Value:          v.Value,
		DefaultWhenOn:  defaultWhenOn,
		DefaultWhenOff: defaultWhenOff,
	}
	flagID, err := primitive.ObjectIDFromHex(flagIDHex)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": flagID}
	res, err := r.flagRepo.col.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{"variants": vrntModel},
	})
	if err != nil {
		return nil, err
	}
	if res.ModifiedCount == 0 {
		return nil, errors.NotFound("flag")
	}
	return vrntModel.asVariant(), nil
}

// func (r VariantRepository) Update(ctx context.Context, flagID string, v variant.UpdateVariantDto) (*flaggio.Variant, error) {
// 	return nil, errors.ErrNotImplemented
// }

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
	})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.NotFound("variant")
	}
	return nil
}

func NewMongoVariantRepository(flagRepo *FlagRepository) *VariantRepository {
	return &VariantRepository{
		flagRepo: flagRepo,
	}
}
