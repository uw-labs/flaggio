package mongodb

import (
	"time"

	"github.com/victorkt/flaggio/internal/flaggio"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type flagModel struct {
	ID                    primitive.ObjectID `bson:"_id"`
	Key                   string             `bson:"key"`
	Name                  string             `bson:"name"`
	Description           *string            `bson:"description"`
	Enabled               bool               `bson:"enabled"`
	Version               int                `bson:"version"`
	Variants              []variantModel     `bson:"variants"`
	Rules                 []flagRuleModel    `bson:"rules"`
	DefaultVariantWhenOn  primitive.ObjectID `bson:"defaultVariantWhenOn"`
	DefaultVariantWhenOff primitive.ObjectID `bson:"defaultVariantWhenOff"`
	CreatedAt             time.Time          `bson:"createdAt"`
	UpdatedAt             *time.Time         `bson:"updatedAt"`
}

func (f *flagModel) asFlag() *flaggio.Flag {
	variants := make([]*flaggio.Variant, len(f.Variants))
	variantsMap := make(map[string]*flaggio.Variant, len(f.Variants))
	for idx, vrntModel := range f.Variants {
		vrnt := vrntModel.asVariant()
		variants[idx] = vrnt
		variantsMap[vrnt.ID] = vrnt
	}
	rules := make([]*flaggio.FlagRule, len(f.Rules))
	for idx, rl := range f.Rules {
		rules[idx] = rl.asRule(variantsMap)
	}
	return &flaggio.Flag{
		ID:                    f.ID.Hex(),
		Key:                   f.Key,
		Name:                  f.Name,
		Description:           f.Description,
		Enabled:               f.Enabled,
		Version:               f.Version,
		Variants:              variants,
		Rules:                 rules,
		DefaultVariantWhenOn:  variantsMap[f.DefaultVariantWhenOn.Hex()],
		DefaultVariantWhenOff: variantsMap[f.DefaultVariantWhenOff.Hex()],
		CreatedAt:             f.CreatedAt,
		UpdatedAt:             f.UpdatedAt,
	}
}

type variantModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Description *string            `bson:"description"`
	Value       interface{}        `bson:"value"`
}

func (v variantModel) asVariant() *flaggio.Variant {
	return &flaggio.Variant{
		ID:          v.ID.Hex(),
		Description: v.Description,
		Value:       v.Value,
	}
}

type flagRuleModel struct {
	ID            primitive.ObjectID  `bson:"_id"`
	Constraints   []constraintModel   `bson:"constraints"`
	Distributions []distributionModel `bson:"distributions"`
}

func (r flagRuleModel) asRule(vrnts map[string]*flaggio.Variant) *flaggio.FlagRule {
	constraints := make([]*flaggio.Constraint, len(r.Constraints))
	for idx, cnstrnt := range r.Constraints {
		constraints[idx] = cnstrnt.asConstraint()
	}
	distributions := make([]*flaggio.Distribution, len(r.Distributions))
	for idx, dstrbtn := range r.Distributions {
		distributions[idx] = dstrbtn.asDistribution(vrnts)
	}
	return &flaggio.FlagRule{
		Rule: flaggio.Rule{
			ID:          r.ID.Hex(),
			Constraints: constraints,
		},
		Distributions: distributions,
	}
}

type constraintModel struct {
	ID        primitive.ObjectID `bson:"_id"`
	Property  string             `bson:"property"`
	Operation string             `bson:"operation"`
	Values    []interface{}      `bson:"values"`
}

func (c constraintModel) asConstraint() *flaggio.Constraint {
	return &flaggio.Constraint{
		ID:        c.ID.Hex(),
		Property:  c.Property,
		Operation: flaggio.Operation(c.Operation),
		Values:    c.Values,
	}
}

type distributionModel struct {
	ID         primitive.ObjectID `bson:"_id"`
	VariantID  primitive.ObjectID `bson:"variantId"`
	Percentage int                `bson:"percentage"`
}

func (d distributionModel) asDistribution(vrnts map[string]*flaggio.Variant) *flaggio.Distribution {
	return &flaggio.Distribution{
		ID:         d.ID.Hex(),
		Variant:    vrnts[d.VariantID.Hex()],
		Percentage: d.Percentage,
	}
}

type segmentRuleModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Constraints []constraintModel  `bson:"constraints"`
}

func (r segmentRuleModel) asRule() *flaggio.SegmentRule {
	constraints := make([]*flaggio.Constraint, len(r.Constraints))
	for idx, cnstrnt := range r.Constraints {
		constraints[idx] = cnstrnt.asConstraint()
	}
	return &flaggio.SegmentRule{
		Rule: flaggio.Rule{
			ID:          r.ID.Hex(),
			Constraints: constraints,
		},
	}
}

type segmentModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description *string            `bson:"description"`
	Rules       []segmentRuleModel `bson:"rules"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   *time.Time         `bson:"updatedAt"`
}

func (f *segmentModel) asSegment() *flaggio.Segment {
	rules := make([]*flaggio.SegmentRule, len(f.Rules))
	for idx, rl := range f.Rules {
		rules[idx] = rl.asRule()
	}
	return &flaggio.Segment{
		ID:          f.ID.Hex(),
		Name:        f.Name,
		Description: f.Description,
		Rules:       rules,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

type evaluationModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	FlagID      primitive.ObjectID `bson:"flagId"`
	FlagKey     string             `bson:"flagKey"`
	FlagVersion int                `bson:"flagVersion"`
	RequestHash string             `bson:"requestHash"`
	UserID      string             `bson:"userId"`
	Value       interface{}        `bson:"value"`
	CreatedAt   time.Time          `bson:"createdAt"`
}

func (f *evaluationModel) asEvaluation() *flaggio.Evaluation {
	return &flaggio.Evaluation{
		ID:          f.ID.Hex(),
		CreatedAt:   f.CreatedAt,
		FlagID:      f.FlagID.Hex(),
		FlagKey:     f.FlagKey,
		FlagVersion: f.FlagVersion,
		RequestHash: f.RequestHash,
		Value:       f.Value,
	}
}

type userModel struct {
	ID        string                 `bson:"_id"`
	Context   map[string]interface{} `bson:"context"`
	UpdatedAt time.Time              `bson:"updatedAt"`
}

func (f *userModel) asUser() *flaggio.User {
	return &flaggio.User{
		ID:        f.ID,
		Context:   f.Context,
		UpdatedAt: f.UpdatedAt,
	}
}
