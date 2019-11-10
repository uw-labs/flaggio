package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

func TestEndsWith_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "ends with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"ef"},
			expectedResult: true,
		},
		{
			desc:           "ends with any strings from the list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "def"},
			expectedResult: true,
		},
		{
			desc:           "doesnt end with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"abc"},
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
			res, err := operator.EndsWith(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestDoesntEndWith_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "doesnt end with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"ab"},
			expectedResult: true,
		},
		{
			desc:           "doesnt end with any strings from list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "abc"},
			expectedResult: true,
		},
		{
			desc:           "ends with string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"abc", "def"},
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
			res, err := operator.DoesntEndWith(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
