package mongodb

import (
	"context"
	"time"

	"github.com/victorkohl/flaggio/internal/errors"
	"github.com/victorkohl/flaggio/internal/flaggio"
	"github.com/victorkohl/flaggio/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ repository.Rule = RuleRepository{}

type RuleRepository struct {
	flagRepo *FlagRepository
}

func (r RuleRepository) CreateFlagRule(ctx context.Context, flagIDHex string, fr flaggio.NewFlagRule) (string, error) {
	constraints := make([]constraintModel, len(fr.Constraints))
	distributions := make([]distributionModel, len(fr.Distributions))
	for idx, c := range fr.Constraints {
		constraints[idx] = constraintModel{
			ID:        primitive.NewObjectID(),
			Property:  c.Property,
			Operation: string(c.Operation),
			Values:    c.Values,
		}
	}
	for idx, d := range fr.Distributions {
		variantID, err := primitive.ObjectIDFromHex(d.VariantID)
		if err != nil {
			return "", errors.New("invalid variant ID for distribution[%d]", idx)
		}
		distributions[idx] = distributionModel{
			ID:         primitive.NewObjectID(),
			VariantID:  variantID,
			Percentage: d.Percentage,
		}
	}
	flgRuleModel := &flagRuleModel{
		ID:            primitive.NewObjectID(),
		Constraints:   constraints,
		Distributions: distributions,
	}
	flagID, err := primitive.ObjectIDFromHex(flagIDHex)
	if err != nil {
		return "", err
	}
	// TODO: check if passing flag id is needed, add unique indexes on all _id fields
	filter := bson.M{"_id": flagID}
	res, err := r.flagRepo.col.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{"rules": flgRuleModel},
		"$set":  bson.M{"updatedAt": time.Now()},
		"$inc":  bson.M{"version": 1},
	})
	if err != nil {
		return "", err
	}
	if res.ModifiedCount == 0 {
		return "", errors.NotFound("flag")
	}
	return flgRuleModel.ID.Hex(), nil
}

func (r RuleRepository) UpdateFlagRule(ctx context.Context, flagIDHex, idHex string, fr flaggio.UpdateFlagRule) error {
	flagID, err := primitive.ObjectIDFromHex(flagIDHex)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	constraints := make([]constraintModel, len(fr.Constraints))
	distributions := make([]distributionModel, len(fr.Distributions))
	for idx, c := range fr.Constraints {
		constraints[idx] = constraintModel{
			ID:        primitive.NewObjectID(),
			Property:  c.Property,
			Operation: string(c.Operation),
			Values:    c.Values,
		}
	}
	for idx, d := range fr.Distributions {
		variantID, err := primitive.ObjectIDFromHex(d.VariantID)
		if err != nil {
			return errors.New("invalid variant ID for distribution[%d]", idx)
		}
		distributions[idx] = distributionModel{
			ID:         primitive.NewObjectID(),
			VariantID:  variantID,
			Percentage: d.Percentage,
		}
	}
	mods := bson.M{
		"updatedAt":             time.Now(),
		"rules.$.constraints":   constraints,
		"rules.$.distributions": distributions,
	}
	res, err := r.flagRepo.col.UpdateOne(
		ctx,
		bson.M{"_id": flagID, "rules._id": id},
		bson.M{"$set": mods, "$inc": bson.M{"version": 1}},
	)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.NotFound("flag rule")
	}
	return nil
}

func (r RuleRepository) DeleteFlagRule(ctx context.Context, flagIDHex, idHex string) error {
	flagID, err := primitive.ObjectIDFromHex(flagIDHex)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	res, err := r.flagRepo.col.UpdateOne(ctx, bson.M{"_id": flagID}, bson.M{
		"$pull": bson.M{"rules": bson.M{"_id": id}},
		"$set":  bson.M{"updatedAt": time.Now()},
		"$inc":  bson.M{"version": 1},
	})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.NotFound("flag rule")
	}
	return nil
}

func NewMongoRuleRepository(flagRepo *FlagRepository) *RuleRepository {
	return &RuleRepository{
		flagRepo: flagRepo,
	}
}
