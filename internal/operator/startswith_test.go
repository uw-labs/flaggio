package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uw-labs/flaggio/internal/operator"
)

func TestStartsWith(t *testing.T) {
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "starts with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"ab"},
			expectedResult: true,
		},
		{
			name:           "starts with any strings from the list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "abc"},
			expectedResult: true,
		},
		{
			name:           "doesnt start with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "unknown config type",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
		{
			name:           "non-string user type",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res, err := operator.StartsWith(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestDoesntStartWith(t *testing.T) {
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "doesnt start with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"cd"},
			expectedResult: true,
		},
		{
			name:           "doesnt start with any strings from list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "cde"},
			expectedResult: true,
		},
		{
			name:           "starts with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"abc", "cde"},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "unknown config type",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: true,
		},
		{
			name:           "non-string user type",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res, err := operator.DoesntStartWith(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
