package flaggio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/flaggio/internal/flaggio"
)

func TestRule_GetID(t *testing.T) {
	t.Parallel()
	rl := flaggio.Rule{ID: "123456"}
	assert.Equal(t, "123456", rl.GetID())
}

func TestFlagRule_Evaluate(t *testing.T) {
	t.Parallel()
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

	tests := []struct {
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

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			eval, err := tt.rule.Evaluate(tt.usrContext)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, eval)
		})
	}
}
