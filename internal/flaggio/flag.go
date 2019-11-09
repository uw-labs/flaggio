package flaggio

import (
	"time"

	"github.com/victorkt/flaggio/internal/errors"
)

var _ Identifier = (*Flag)(nil)
var _ Evaluator = (*Flag)(nil)

type Flag struct {
	ID                    string
	Key                   string
	Name                  string
	Description           *string
	Enabled               bool
	Version               int
	Variants              []*Variant
	Rules                 []*FlagRule
	DefaultVariantWhenOn  *Variant
	DefaultVariantWhenOff *Variant
	CreatedAt             time.Time
	UpdatedAt             *time.Time
}

func (f Flag) GetID() string {
	return f.ID
}

func (f Flag) Evaluate(usrContext map[string]interface{}) (EvalResult, error) {
	var answer interface{}
	var next []Evaluator
	if f.Enabled {
		vrnt := f.DefaultVariantWhenOn
		if vrnt == nil {
			return EvalResult{}, errors.ErrNoDefaultVariant
		}
		answer = vrnt.Value
		for _, rl := range f.Rules {
			next = append(next, rl)
		}
	} else {
		vrnt := f.DefaultVariantWhenOff
		if vrnt == nil {
			return EvalResult{}, errors.ErrNoDefaultVariant
		}
		answer = vrnt.Value
	}
	return EvalResult{
		Answer: answer,
		Next:   next,
	}, nil
}

func (f *Flag) Populate(identifiers []Identifier) {
	for _, r := range f.Rules {
		r.Populate(identifiers)
	}
}
