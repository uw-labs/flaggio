package flaggio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/flaggio/internal/flaggio"
)

func TestSegment_GetID(t *testing.T) {
	t.Parallel()
	sgmnt := flaggio.Segment{ID: "123456"}
	assert.Equal(t, "123456", sgmnt.GetID())
}

func TestSegment_Validate(t *testing.T) {
	t.Parallel()
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

	tests := []struct {
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

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.segment.Validate(tt.usrContext)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
