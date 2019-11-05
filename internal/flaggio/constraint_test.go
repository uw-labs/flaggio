package flaggio

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	flaggio_mock "github.com/victorkt/flaggio/internal/flaggio/mocks"
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
	}{
		{
			desc:           "returns false on unknown operators",
			cnstrnt:        Constraint{Property: "key", Operation: "INVALID", Values: []interface{}{}},
			operatorCalls:  0,
			expectedResult: false,
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			op := flaggio_mock.NewMockOperator(mockCtrl)
			oldOperatorMap := make(map[Operation]Operator, len(operatorMap))
			for key, val := range operatorMap {
				oldOperatorMap[key] = val
				operatorMap[key] = op
			}
			defer func() {
				operatorMap = oldOperatorMap
			}()

			op.EXPECT().Operate(test.expectedUsrValue, gomock.Any()).Times(test.operatorCalls).Return(test.operatorResult)

			res := test.cnstrnt.Validate(test.usrContext)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
