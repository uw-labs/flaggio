package mongodb

import (
	"time"

	"github.com/victorkohl/flaggio/internal/flaggio"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type flagModel struct {
	ID          primitive.ObjectID `bson:"_id"`
	Key         string             `bson:"key"`
	Name        string             `bson:"name"`
	Description *string            `bson:"description"`
	Enabled     bool               `bson:"enabled"`
	Version     int                `bson:"version"`
	Variants    []variantModel     `bson:"variants"`
	Rules       []flagRuleModel    `bson:"rules"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   *time.Time         `bson:"updatedAt"`
}

func (f flagModel) asFlag() *flaggio.Flag {
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
		ID:          f.ID.Hex(),
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
		Key:         f.Key,
		Name:        f.Name,
		Description: f.Description,
		Enabled:     f.Enabled,
		Version:     f.Version,
		Variants:    variants,
		Rules:       rules,
	}
}

type variantModel struct {
	ID             primitive.ObjectID `bson:"_id"`
	Key            string             `bson:"key"`
	Description    *string            `bson:"description"`
	Value          interface{}        `bson:"value"`
	DefaultWhenOn  bool               `bson:"defaultWhenOn"`
	DefaultWhenOff bool               `bson:"defaultWhenOff"`
}

func (v variantModel) asVariant() *flaggio.Variant {
	return &flaggio.Variant{
		ID:             v.ID.Hex(),
		Key:            v.Key,
		Description:    v.Description,
		Value:          v.Value,
		DefaultWhenOn:  v.DefaultWhenOn,
		DefaultWhenOff: v.DefaultWhenOff,
	}
}

type flagRuleModel struct {
	ID            primitive.ObjectID `bson:"_id"`
	Constraints   []constraintModel  `bson:"constraints"`
	Distributions []distribution     `bson:"distributions"`
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

type distribution struct {
	ID         primitive.ObjectID `bson:"_id"`
	VariantID  primitive.ObjectID `bson:"variantId"`
	Percentage uint32             `bson:"percentage"`
}

func (d distribution) asDistribution(vrnts map[string]*flaggio.Variant) *flaggio.Distribution {
	return &flaggio.Distribution{
		Variant:    vrnts[d.VariantID.Hex()],
		Percentage: int(d.Percentage),
	}
}
