package flaggio

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstraint_Validate(t *testing.T) {
	tests := []struct {
		name             string
		cnstrnt          Constraint
		usrContext       map[string]interface{}
		operatorCalls    int
		operatorResult   bool
		expectedUsrValue interface{}
		expectedResult   bool
		expectedError    error
	}{
		{
			name:           "returns error on unknown operators",
			cnstrnt:        Constraint{Property: "key", Operation: "INVALID", Values: []interface{}{}},
			operatorCalls:  0,
			expectedResult: false,
			expectedError:  errors.New("invalid flag: unknown operation: INVALID"),
		},
		{
			name:             "returns true when operator returns true",
			cnstrnt:          Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{1}},
			usrContext:       map[string]interface{}{"key": 1},
			operatorCalls:    1,
			operatorResult:   true,
			expectedUsrValue: 1,
			expectedResult:   true,
		},
		{
			name:             "returns false when operator returns false",
			cnstrnt:          Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{1}},
			usrContext:       map[string]interface{}{"key": 1},
			operatorCalls:    1,
			operatorResult:   false,
			expectedUsrValue: 1,
			expectedResult:   false,
		},
		{
			name:             "passes the user context as argument for IsInSegment operations",
			cnstrnt:          Constraint{Property: "key", Operation: OperationIsInSegment, Values: []interface{}{1}},
			usrContext:       map[string]interface{}{"key": 1},
			operatorCalls:    1,
			operatorResult:   false,
			expectedUsrValue: map[string]interface{}{"key": 1},
			expectedResult:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var mockCalls int
			mockOp := func(usrValue interface{}, validValues []interface{}) (bool, error) {
				mockCalls++
				if !reflect.DeepEqual(usrValue, tt.expectedUsrValue) {
					return false, fmt.Errorf("expected user value to be: %+v, got: %+v", tt.expectedUsrValue, usrValue)
				}
				return tt.operatorResult, nil
			}
			oldOperator, ok := operatorMap[tt.cnstrnt.Operation]
			if ok {
				operatorMap[tt.cnstrnt.Operation] = mockOp
				defer func() {
					operatorMap[tt.cnstrnt.Operation] = oldOperator
				}()
			}

			res, err := tt.cnstrnt.Validate(tt.usrContext)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.operatorCalls, mockCalls)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestConstraintList_Validate(t *testing.T) {
	tests := []struct {
		name           string
		cnstrnts       ConstraintList
		usrContext     map[string]interface{}
		operatorResult bool
		expectedResult bool
	}{
		{
			name: "returns false when any constraint is false",
			cnstrnts: ConstraintList{
				&Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{1}},
				&Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{2}},
			},
			usrContext:     map[string]interface{}{"key": 1},
			expectedResult: false,
		},
		{
			name: "returns true when all constraints returns true",
			cnstrnts: ConstraintList{
				&Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{1}},
				&Constraint{Property: "key2", Operation: OperationOneOf, Values: []interface{}{2}},
			},
			usrContext:     map[string]interface{}{"key": 1, "key2": 2},
			operatorResult: true,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.cnstrnts.Validate(tt.usrContext)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
