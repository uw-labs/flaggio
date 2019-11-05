package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

var _ operator.Validator = mockValidator{}

type mockValidator struct {
	r bool
}

func (m mockValidator) Validate(usrContext map[string]interface{}) bool {
	return m.r
}

func TestValidates_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "validates",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true}},
			expectedResult: true,
		},
		{
			desc:           "validates all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true}, mockValidator{true}},
			expectedResult: true,
		},
		{
			desc:           "doesnt validate",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false}},
			expectedResult: false,
		},
		{
			desc:           "doesnt validate all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true}, mockValidator{false}},
			expectedResult: false,
		},
		{
			desc:           "doesnt validate when not a validator",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			op := operator.Validates{}
			res := op.Operate(test.usrContext, test.values)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestDoesntValidate_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "validates",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false}},
			expectedResult: true,
		},
		{
			desc:           "validates all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false}, mockValidator{false}},
			expectedResult: true,
		},
		{
			desc:           "doesnt validate",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true}},
			expectedResult: false,
		},
		{
			desc:           "doesnt validate all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false}, mockValidator{true}},
			expectedResult: false,
		},
		{
			desc:           "doesnt validate when not a validator",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			op := operator.DoesntValidate{}
			res := op.Operate(test.usrContext, test.values)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
