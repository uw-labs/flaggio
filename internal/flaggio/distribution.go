package flaggio

import (
	"math/rand"
	"time"

	"github.com/uw-labs/flaggio/internal/errors"
)

var _ Evaluator = (*DistributionList)(nil)

// Distribution represents a percentage chance for a variant to
// be selected as the result value for the flag evaluation.
type Distribution struct {
	ID         string
	Variant    *Variant
	Percentage int
}

// DistributionList is a slice of *Distribution.
type DistributionList []*Distribution

// Evaluate will select one of the distributions based on their percentage
// chance and return it's value as answer.
func (dl DistributionList) Evaluate(usrContext map[string]interface{}) (EvalResult, error) {
	ref := dl.Distribute()
	if ref == nil || ref.ID == "" {
		// configuration problem, return error
		return EvalResult{}, errors.ErrNoVariantToDistribute
	}
	return EvalResult{
		Answer: ref.Value,
	}, nil
}

// Distribute selects a distribution randomly, respecting the configured probability.
func (dl DistributionList) Distribute() *Variant {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint:gosec // not security critical
	num := 1 + r1.Intn(100)                               // random int between 1 and 100

	var total int
	for _, dstrbtn := range dl {
		total += dstrbtn.Percentage
		if num <= total {
			return dstrbtn.Variant
		}
	}

	// fallback, should never happen
	return dl[0].Variant
}
