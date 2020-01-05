package flaggio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
)

func TestRule_GetID(t *testing.T) {
	rl := flaggio.Rule{ID: "123456"}
	assert.Equal(t, "123456", rl.GetID())
}

func TestFlagRule_Evaluate(t *testing.T) {
	vrnt1 := &flaggio.Variant{ID: "1", Value: 1}
	dstrbtn := &flaggio.Distribution{
		ID:         "1",
		Variant:    vrnt1,
		Percentage: 100,
	}
	rl := flaggio.Rule{
		ID: "123-abc",
		Constraints: []*flaggio.Constraint{{
			ID:        "1",
			Property:  "name",
			Operation: flaggio.OperationOneOf,
			Values:    []interface{}{"John"},
		}, {
			ID:        "2",
			Property:  "age",
			Operation: flaggio.OperationGreaterOrEqual,
			Values:    []interface{}{30},
		}},
	}

	tt := []struct {
		name           string
		usrContext     map[string]interface{}
		rule           flaggio.FlagRule
		expectedResult flaggio.EvalResult
		expectedError  error
	}{
		{
			name:       "doesn't return the distribution list when ANY constraint validate to false",
			usrContext: map[string]interface{}{"name": "Mary", "age": 40},
			rule: flaggio.FlagRule{
				Rule:          rl,
				Distributions: []*flaggio.Distribution{dstrbtn},
			},
			expectedResult: flaggio.EvalResult{Answer: nil, Next: nil},
		},
		{
			name:       "returns the distribution list when ALL constraints validate to true",
			usrContext: map[string]interface{}{"name": "John", "age": 30},
			rule: flaggio.FlagRule{
				Rule:          rl,
				Distributions: []*flaggio.Distribution{dstrbtn},
			},
			expectedResult: flaggio.EvalResult{
				Answer: nil,
				Next:   []flaggio.Evaluator{flaggio.DistributionList([]*flaggio.Distribution{dstrbtn})},
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			eval, err := test.rule.Evaluate(test.usrContext)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedResult, eval)
		})
	}
}
