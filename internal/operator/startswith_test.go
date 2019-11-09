package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

func TestStartsWith_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "starts with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"ab"},
			expectedResult: true,
		},
		{
			desc:           "starts with any strings from the list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "abc"},
			expectedResult: true,
		},
		{
			desc:           "doesnt start with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			op := operator.StartsWith{}
			res, err := op.Operate(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestDoesntStartWith_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "doesnt start with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"cd"},
			expectedResult: true,
		},
		{
			desc:           "doesnt start with any strings from list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "cde"},
			expectedResult: true,
		},
		{
			desc:           "starts with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"abc", "cde"},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: true,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			op := operator.DoesntStartWith{}
			res, err := op.Operate(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
