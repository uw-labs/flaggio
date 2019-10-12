package flaggio

import (
	"time"

	"github.com/victorkohl/flaggio/internal/errors"
)

type Flag struct {
	ID          string
	Key         string
	Name        string
	Description *string
	Enabled     bool
	Version     uint64
	Variants    []*Variant
	Rules       []*FlagRule
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func (f Flag) Evaluate(usr map[string]interface{}) (EvalResult, error) {
	var answer interface{}
	var next Evaluator
	if f.Enabled {
		vrnt := VariantList(f.Variants).DefaultWhenOn()
		if vrnt == nil {
			return EvalResult{}, errors.ErrNoDefaultVariant
		}
		answer = vrnt.Value
		// next = f.Rules
	} else {
		vrnt := VariantList(f.Variants).DefaultWhenOff()
		if vrnt == nil {
			return EvalResult{}, errors.ErrNoDefaultVariant
		}
		answer = vrnt.Value
	}
	return EvalResult{
		Answer: answer,
		Next:   []Evaluator{next},
	}, nil
}
