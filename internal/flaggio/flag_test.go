package flaggio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/errors"
	"github.com/victorkt/flaggio/internal/flaggio"
)

func TestFlag_GetID(t *testing.T) {
	flg := flaggio.Flag{ID: "123456"}
	assert.Equal(t, "123456", flg.GetID())
}

func TestFlag_Evaluate(t *testing.T) {
	vrnt1 := &flaggio.Variant{ID: "1", Value: 1}
	vrnt2 := &flaggio.Variant{ID: "2", Value: 2}
	rl1 := &flaggio.FlagRule{}

	tt := []struct {
		name           string
		flag           flaggio.Flag
		expectedResult flaggio.EvalResult
		expectedError  error
	}{
		{
			name: "returns error when there is no variant when off",
			flag: flaggio.Flag{
				Enabled:  false,
				Variants: []*flaggio.Variant{vrnt1, vrnt2},
			},
			expectedError: errors.ErrNoDefaultVariant,
		},
		{
			name: "returns error when there is no variant when on",
			flag: flaggio.Flag{
				Enabled:  true,
				Variants: []*flaggio.Variant{vrnt1, vrnt2},
			},
			expectedError: errors.ErrNoDefaultVariant,
		},
		{
			name: "returns default variant when off",
			flag: flaggio.Flag{
				Enabled:               false,
				Variants:              []*flaggio.Variant{vrnt1, vrnt2},
				DefaultVariantWhenOn:  vrnt1,
				DefaultVariantWhenOff: vrnt2,
			},
			expectedResult: flaggio.EvalResult{Answer: 2},
		},
		{
			name: "returns default variant when on",
			flag: flaggio.Flag{
				Enabled:               true,
				Variants:              []*flaggio.Variant{vrnt1, vrnt2},
				Rules:                 []*flaggio.FlagRule{rl1},
				DefaultVariantWhenOn:  vrnt1,
				DefaultVariantWhenOff: vrnt2,
			},
			expectedResult: flaggio.EvalResult{Answer: 1, Next: []flaggio.Evaluator{rl1}},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			usrContext := map[string]interface{}{"name": "John"}
			eval, err := test.flag.Evaluate(usrContext)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedResult, eval)
		})
	}
}
