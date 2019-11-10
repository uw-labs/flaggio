package flaggio

import (
	"time"

	"github.com/victorkt/flaggio/internal/errors"
)

var _ Identifier = (*Flag)(nil)
var _ Evaluator = (*Flag)(nil)

// Flag holds all the information needed to evaluate a key to a value.
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

// GetID returns the flag ID.
func (f Flag) GetID() string {
	return f.ID
}

// Evaluate will return the default variant as answer based on the flag status (on or off).
// If the flag is on, it will also return the list of rules to be evaluated.
// If there is no default variant configured for the given flag enabled state, an error
// is returned.
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

// Populate will try to populate all references in the list of rules.
func (f *Flag) Populate(identifiers []Identifier) {
	for _, r := range f.Rules {
		r.Populate(identifiers)
	}
}
