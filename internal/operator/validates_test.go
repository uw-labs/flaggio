package operator_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/flaggio/internal/operator"
)

var _ operator.Validator = mockValidator{}

type mockValidator struct {
	r bool
	e error
}

func (m mockValidator) Validate(usrContext map[string]interface{}) (bool, error) {
	return m.r, m.e
}

func TestValidates(t *testing.T) {
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "validates",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true, nil}},
			expectedResult: true,
		},
		{
			name:           "validates all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true, nil}, mockValidator{true, nil}},
			expectedResult: true,
		},
		{
			name:           "doesnt validate",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, nil}},
			expectedResult: false,
		},
		{
			name:           "doesnt validate all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true, nil}, mockValidator{false, nil}},
			expectedResult: false,
		},
		{
			name:           "returns error",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, errors.New("test")}},
			expectedResult: false,
			expectedError:  errors.New("test"),
		},
		{
			name:           "doesnt validate when not a validator",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res, err := operator.Validates(tt.usrContext, tt.values)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestDoesntValidate(t *testing.T) {
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "validates",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, nil}},
			expectedResult: true,
		},
		{
			name:           "validates all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, nil}, mockValidator{false, nil}},
			expectedResult: true,
		},
		{
			name:           "doesnt validate",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{true, nil}},
			expectedResult: false,
		},
		{
			name:           "doesnt validate all",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, nil}, mockValidator{true, nil}},
			expectedResult: false,
		},
		{
			name:           "returns error",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{mockValidator{false, errors.New("test")}},
			expectedResult: false,
			expectedError:  errors.New("test"),
		},
		{
			name:           "doesnt validate when not a validator",
			usrContext:     map[string]interface{}{"prop": "abc"},
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res, err := operator.DoesntValidate(tt.usrContext, tt.values)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
