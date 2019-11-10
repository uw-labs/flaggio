package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

func TestContains_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "contains string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: true,
		},
		{
			desc:           "contains any strings from the list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "cde"},
			expectedResult: true,
		},
		{
			desc:           "doesnt contain string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"xyz"},
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
			res, err := operator.Contains(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestDoesntContain_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "doesnt contain string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"xyz"},
			expectedResult: true,
		},
		{
			desc:           "doesnt contain any strings from list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "xyz"},
			expectedResult: true,
		},
		{
			desc:           "contains string",
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
			res, err := operator.DoesntContain(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
