package flaggio

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstraint_Validate(t *testing.T) {
	tt := []struct {
		desc             string
		cnstrnt          Constraint
		usrContext       map[string]interface{}
		operatorCalls    int
		operatorResult   bool
		expectedUsrValue interface{}
		expectedResult   bool
		expectedError    error
	}{
		{
			desc:           "returns error on unknown operators",
			cnstrnt:        Constraint{Property: "key", Operation: "INVALID", Values: []interface{}{}},
			operatorCalls:  0,
			expectedResult: false,
			expectedError:  errors.New("invalid flag: unknown operation: INVALID"),
		},
		{
			desc:             "returns true when operator returns true",
			cnstrnt:          Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{1}},
			usrContext:       map[string]interface{}{"key": 1},
			operatorCalls:    1,
			operatorResult:   true,
			expectedUsrValue: 1,
			expectedResult:   true,
		},
		{
			desc:             "returns false when operator returns false",
			cnstrnt:          Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{1}},
			usrContext:       map[string]interface{}{"key": 1},
			operatorCalls:    1,
			operatorResult:   false,
			expectedUsrValue: 1,
			expectedResult:   false,
		},
		{
			desc:             "passes the user context as argument for IsInSegment operations",
			cnstrnt:          Constraint{Property: "key", Operation: OperationIsInSegment, Values: []interface{}{1}},
			usrContext:       map[string]interface{}{"key": 1},
			operatorCalls:    1,
			operatorResult:   false,
			expectedUsrValue: map[string]interface{}{"key": 1},
			expectedResult:   false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			var mockCalls int
			mockOp := func(usrValue interface{}, validValues []interface{}) (bool, error) {
				mockCalls++
				if !reflect.DeepEqual(usrValue, test.expectedUsrValue) {
					return false, fmt.Errorf("expected user value to be: %+v, got: %+v", test.expectedUsrValue, usrValue)
				}
				return test.operatorResult, nil
			}
			oldOperator, ok := operatorMap[test.cnstrnt.Operation]
			if ok {
				operatorMap[test.cnstrnt.Operation] = mockOp
				defer func() {
					operatorMap[test.cnstrnt.Operation] = oldOperator
				}()
			}

			res, err := test.cnstrnt.Validate(test.usrContext)
			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.operatorCalls, mockCalls)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestConstraintList_Validate(t *testing.T) {
	tt := []struct {
		desc           string
		cnstrnts       ConstraintList
		usrContext     map[string]interface{}
		operatorResult bool
		expectedResult bool
	}{
		{
			desc: "returns false when any constraint is false",
			cnstrnts: ConstraintList{
				&Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{1}},
				&Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{2}},
			},
			usrContext:     map[string]interface{}{"key": 1},
			expectedResult: false,
		},
		{
			desc: "returns true when all constraints returns true",
			cnstrnts: ConstraintList{
				&Constraint{Property: "key", Operation: OperationOneOf, Values: []interface{}{1}},
				&Constraint{Property: "key2", Operation: OperationOneOf, Values: []interface{}{2}},
			},
			usrContext:     map[string]interface{}{"key": 1, "key2": 2},
			operatorResult: true,
			expectedResult: true,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			res, err := test.cnstrnts.Validate(test.usrContext)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
