package flaggio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
)

func TestSegment_GetID(t *testing.T) {
	sgmnt := flaggio.Segment{ID: "123456"}
	assert.Equal(t, "123456", sgmnt.GetID())
}

func TestSegment_Validate(t *testing.T) {
	rl1 := flaggio.Rule{
		Constraints: []*flaggio.Constraint{{
			Property:  "name",
			Operation: flaggio.OperationOneOf,
			Values:    []interface{}{"John"},
		}},
	}
	rl2 := flaggio.Rule{
		Constraints: []*flaggio.Constraint{{
			Property:  "age",
			Operation: flaggio.OperationGreaterOrEqual,
			Values:    []interface{}{30},
		}},
	}

	tt := []struct {
		name           string
		usrContext     map[string]interface{}
		segment        flaggio.Segment
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "returns true when ANY of the rules validate to true",
			usrContext:     map[string]interface{}{"name": "Mary", "age": 40},
			segment:        flaggio.Segment{Rules: []*flaggio.SegmentRule{{rl1}, {rl2}}},
			expectedResult: true,
		},
		{
			name:           "returns false when ALL of the rules validate to false",
			usrContext:     map[string]interface{}{"name": "Jane", "age": 25},
			segment:        flaggio.Segment{Rules: []*flaggio.SegmentRule{{rl1}, {rl2}}},
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.segment.Validate(test.usrContext)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}
