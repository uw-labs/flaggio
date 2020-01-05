package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

func TestContains_Operate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "contains string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: true,
		},
		{
			name:           "contains any strings from the list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "cde"},
			expectedResult: true,
		},
		{
			name:           "doesnt contain string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"xyz"},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := operator.Contains(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestDoesntContain_Operate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "doesnt contain string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"xyz"},
			expectedResult: true,
		},
		{
			name:           "doesnt contain any strings from list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "xyz"},
			expectedResult: true,
		},
		{
			name:           "contains string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"abc", "cde"},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := operator.DoesntContain(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
