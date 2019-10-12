package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkohl/flaggio/internal/operator"
)

func TestExists_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usr            map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "property exists",
			usr:            map[string]interface{}{"prop": "abc"},
			property:       "prop",
			expectedResult: true,
		},
		{
			desc:           "property doesnt exist",
			usr:            map[string]interface{}{"prop": "abc"},
			property:       "prop2",
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			op := operator.Exists{}
			res := op.Operate(test.usr[test.property], test.values)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestDoesntExist_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usr            map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "property exists",
			usr:            map[string]interface{}{"prop": "abc"},
			property:       "prop",
			expectedResult: false,
		},
		{
			desc:           "property doesnt exist",
			usr:            map[string]interface{}{"prop": "abc"},
			property:       "prop2",
			expectedResult: true,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			op := operator.DoesntExist{}
			res := op.Operate(test.usr[test.property], test.values)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
