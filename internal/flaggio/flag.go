package flaggio

import (
	"time"

	"github.com/victorkohl/flaggio/internal/errors"
)

var _ Identifier = Flag{}
var _ Evaluator = Flag{}

type Flag struct {
	ID          string
	Key         string
	Name        string
	Description *string
	Enabled     bool
	Version     int
	Variants    []*Variant
	Rules       []*FlagRule
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func (f Flag) GetID() string {
	return f.ID
}

func (f Flag) Evaluate(usrContext map[string]interface{}) (EvalResult, error) {
	var answer interface{}
	var next []Evaluator
	if f.Enabled {
		vrnt := VariantList(f.Variants).DefaultWhenOn()
		if vrnt == nil {
			return EvalResult{}, errors.ErrNoDefaultVariant
		}
		answer = vrnt.Value
		for _, rl := range f.Rules {
			next = append(next, rl)
		}
	} else {
		vrnt := VariantList(f.Variants).DefaultWhenOff()
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
