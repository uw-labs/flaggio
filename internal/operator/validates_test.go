package operator_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

var _ operator.Validator = mockValidator{}

type mockValidator struct {
	r bool
	e error
}

func (m mockValidator) Validate(usrContext map[string]interface{}) (bool, error) {
	return m.r, m.e
}

func TestValidates_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
		expectedError  error
	}{
		{
			desc:           "validates",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true, nil}},
			expectedResult: true,
		},
		{
			desc:           "validates all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true, nil}, mockValidator{true, nil}},
			expectedResult: true,
		},
		{
			desc:           "doesnt validate",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, nil}},
			expectedResult: false,
		},
		{
			desc:           "doesnt validate all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true, nil}, mockValidator{false, nil}},
			expectedResult: false,
		},
		{
			desc:           "returns error",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, errors.New("test")}},
			expectedResult: false,
			expectedError:  errors.New("test"),
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
			res, err := operator.Validates(test.usrContext, test.values)
			assert.Equal(t, test.expectedError, err)
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
		expectedError  error
	}{
		{
			desc:           "validates",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, nil}},
			expectedResult: true,
		},
		{
			desc:           "validates all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, nil}, mockValidator{false, nil}},
			expectedResult: true,
		},
		{
			desc:           "doesnt validate",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true, nil}},
			expectedResult: false,
		},
		{
			desc:           "doesnt validate all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, nil}, mockValidator{true, nil}},
			expectedResult: false,
		},
		{
			desc:           "returns error",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, errors.New("test")}},
			expectedResult: false,
			expectedError:  errors.New("test"),
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
			res, err := operator.DoesntValidate(test.usrContext, test.values)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
