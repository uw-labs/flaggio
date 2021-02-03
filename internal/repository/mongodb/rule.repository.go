package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uw-labs/flaggio/internal/errors"
	"github.com/uw-labs/flaggio/internal/flaggio"
	"github.com/uw-labs/flaggio/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ repository.Rule = (*RuleRepository)(nil)

// RuleRepository implements repository.Rule interface using mongodb.
type RuleRepository struct {
	flagRepo    *FlagRepository
	segmentRepo *SegmentRepository
}

// FindFlagRuleByID returns a flag rule that has a given ID.
func (r *RuleRepository) FindFlagRuleByID(ctx context.Context, flagIDHex, idHex string) (*flaggio.FlagRule, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoRuleRepository.FindFlagRuleByID")
	defer span.Finish()

	flagID, err := primitive.ObjectIDFromHex(flagIDHex)
	if err != nil {
		return nil, err
	}
	ruleID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": flagID, "rules._id": ruleID}
	projection := bson.M{"variants": 1, "rules.$": 1}
	opts := options.FindOne().SetProjection(projection)

	var f flagModel
	if err := r.flagRepo.col.FindOne(ctx, filter, opts).Decode(&f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("rule")
		}
		return nil, err
	}
	if len(f.Rules) != 1 {
		return nil, errors.NotFound("rule")
	}
	flg := f.asFlag()
	return flg.Rules[0], nil
}

// CreateFlagRule creates a new rule under a flag.
func (r *RuleRepository) CreateFlagRule(ctx context.Context, flagIDHex string, fr flaggio.NewFlagRule) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoRuleRepository.CreateFlagRule")
	defer span.Finish()

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
			return "", errors.BadRequest(fmt.Sprintf("invalid variant ID for distribution[%d]", idx))
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

// UpdateFlagRule updates a rule under a flag.
func (r *RuleRepository) UpdateFlagRule(ctx context.Context, flagIDHex, idHex string, fr flaggio.UpdateFlagRule) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoRuleRepository.UpdateFlagRule")
	defer span.Finish()

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
			return errors.BadRequest(fmt.Sprintf("invalid variant ID for distribution[%d]", idx))
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

// DeleteFlagRule deletes a rule under a flag.
func (r *RuleRepository) DeleteFlagRule(ctx context.Context, flagIDHex, idHex string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoRuleRepository.DeleteFlagRule")
	defer span.Finish()

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

// FindSegmentRuleByID returns a segment rule that has a given ID.
func (r *RuleRepository) FindSegmentRuleByID(ctx context.Context, segmentIDHex, idHex string) (*flaggio.SegmentRule, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoRuleRepository.FindSegmentRuleByID")
	defer span.Finish()

	segmentID, err := primitive.ObjectIDFromHex(segmentIDHex)
	if err != nil {
		return nil, err
	}
	ruleID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": segmentID, "rules._id": ruleID}
	projection := bson.M{"rules.$": 1}
	opts := options.FindOne().SetProjection(projection)

	var s segmentModel
	if err := r.segmentRepo.col.FindOne(ctx, filter, opts).Decode(&s); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("rule")
		}
		return nil, err
	}
	if len(s.Rules) != 1 {
		return nil, errors.NotFound("rule")
	}
	return s.Rules[0].asRule(), nil
}

// CreateSegmentRule creates a new rule under a segment.
func (r *RuleRepository) CreateSegmentRule(ctx context.Context, segmentIDHex string, fr flaggio.NewSegmentRule) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoRuleRepository.CreateSegmentRule")
	defer span.Finish()

	constraints := make([]constraintModel, len(fr.Constraints))
	for idx, c := range fr.Constraints {
		constraints[idx] = constraintModel{
			ID:        primitive.NewObjectID(),
			Property:  c.Property,
			Operation: string(c.Operation),
			Values:    c.Values,
		}
	}
	sgmntRuleModel := &segmentRuleModel{
		ID:          primitive.NewObjectID(),
		Constraints: constraints,
	}
	segmentID, err := primitive.ObjectIDFromHex(segmentIDHex)
	if err != nil {
		return "", err
	}
	filter := bson.M{"_id": segmentID}
	res, err := r.segmentRepo.col.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{"rules": sgmntRuleModel},
		"$set":  bson.M{"updatedAt": time.Now()},
	})
	if err != nil {
		return "", err
	}
	if res.ModifiedCount == 0 {
		return "", errors.NotFound("segment")
	}
	return sgmntRuleModel.ID.Hex(), nil
}

// UpdateSegmentRule updates a rule under a segment.
func (r *RuleRepository) UpdateSegmentRule(ctx context.Context, segmentIDHex, idHex string, fr flaggio.UpdateSegmentRule) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoRuleRepository.UpdateSegmentRule")
	defer span.Finish()

	segmentID, err := primitive.ObjectIDFromHex(segmentIDHex)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	constraints := make([]constraintModel, len(fr.Constraints))
	for idx, c := range fr.Constraints {
		constraints[idx] = constraintModel{
			ID:        primitive.NewObjectID(),
			Property:  c.Property,
			Operation: string(c.Operation),
			Values:    c.Values,
		}
	}
	mods := bson.M{
		"updatedAt":           time.Now(),
		"rules.$.constraints": constraints,
	}
	res, err := r.segmentRepo.col.UpdateOne(
		ctx,
		bson.M{"_id": segmentID, "rules._id": id},
		bson.M{"$set": mods},
	)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.NotFound("segment rule")
	}
	return nil
}

// DeleteSegmentRule deletes a rule under a segment.
func (r *RuleRepository) DeleteSegmentRule(ctx context.Context, segmentIDHex, idHex string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoRuleRepository.DeleteSegmentRule")
	defer span.Finish()

	segmentID, err := primitive.ObjectIDFromHex(segmentIDHex)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	res, err := r.segmentRepo.col.UpdateOne(ctx, bson.M{"_id": segmentID}, bson.M{
		"$pull": bson.M{"rules": bson.M{"_id": id}},
		"$set":  bson.M{"updatedAt": time.Now()},
	})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.NotFound("segment rule")
	}
	return nil
}

// NewRuleRepository returns a new rule repository that uses mongodb as underlying storage.
func NewRuleRepository(flagRepo *FlagRepository, segmentRepo *SegmentRepository) repository.Rule {
	return &RuleRepository{
		flagRepo:    flagRepo,
		segmentRepo: segmentRepo,
	}
}
