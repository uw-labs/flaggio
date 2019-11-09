package flaggio

import (
	"math/rand"
	"time"

	"github.com/victorkt/flaggio/internal/errors"
)

var _ Evaluator = (*DistributionList)(nil)

type Distribution struct {
	ID         string
	Variant    *Variant
	Percentage int
}

type DistributionList []*Distribution

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

func (dl DistributionList) Distribute() *Variant {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := 1 + r1.Intn(100) // random num between 1 and 100

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
