package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

func TestExists_Operate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "property exists",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			expectedResult: true,
		},
		{
			name:           "property doesnt exist",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop2",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := operator.Exists(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestDoesntExist_Operate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "property exists",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			expectedResult: false,
		},
		{
			name:           "property doesnt exist",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop2",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := operator.DoesntExist(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
