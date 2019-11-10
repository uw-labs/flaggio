package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

func TestMatchesRegex_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "matches string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"[a-z]+"},
			expectedResult: true,
		},
		{
			desc:           "matches any strings from the list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"[a-z]+", "[0-9]+"},
			expectedResult: true,
		},
		{
			desc:           "doesnt match string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"[0-9]+"},
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
			res, err := operator.MatchesRegex(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestDoesntMatchRegex_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "doesnt match string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"[0-9]+"},
			expectedResult: true,
		},
		{
			desc:           "doesnt match any strings from list",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"[0-9]+", "[x-z]+"},
			expectedResult: true,
		},
		{
			desc:           "matches string",
			usrContext:     map[string]interface{}{"prop": "abcdef"},
			property:       "prop",
			values:         []interface{}{"[a-z]+"},
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
			res, err := operator.DoesntMatchRegex(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
