package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

func TestGreater_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "int greater than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: true,
		},
		{
			desc:           "int greater than int from list",
			usrContext:     map[string]interface{}{"prop": int(2)},
			property:       "prop",
			values:         []interface{}{int(0), int(1)},
			expectedResult: true,
		},
		{
			desc:           "int not greater than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int greater than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(0)},
			expectedResult: true,
		},
		{
			desc:           "int not greater than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			desc:           "int greater than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(0)},
			expectedResult: true,
		},
		{
			desc:           "int not greater than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int32 greater than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int32 greater than int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int32(0), int32(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not greater than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 greater than int",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not greater than int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 greater than int64",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not greater than int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			desc:           "int64 greater than int64",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 greater than int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int64(0), int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 not greater than int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			desc:           "int64 greater than int",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not greater than int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int64 greater than int32",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not greater than int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			desc:           "float64 greater than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.2)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 greater than float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(9.2)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 not greater than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			res, err := operator.Greater(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestGreaterOrEqual_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "int greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			desc:           "int greater or equal than int from list",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0), int(1)},
			expectedResult: true,
		},
		{
			desc:           "int not greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(0)},
			expectedResult: true,
		},
		{
			desc:           "int not greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			desc:           "int greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(0)},
			expectedResult: true,
		},
		{
			desc:           "int not greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int32 greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int32 greater or equal than int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int32(0), int32(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			desc:           "int64 greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 greater or equal than int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int64(0), int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 not greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			desc:           "int64 greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int64 greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			desc:           "float64 greater or equal than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.2)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 greater or equal than float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(9.2)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 not greater or equal than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			res, err := operator.GreaterOrEqual(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestLower_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "int lower than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int lower than int from list",
			usrContext:     map[string]interface{}{"prop": int(2)},
			property:       "prop",
			values:         []interface{}{int(3), int(4)},
			expectedResult: true,
		},
		{
			desc:           "int not lower than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: false,
		},
		{
			desc:           "int lower than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		{
			desc:           "int not lower than int32",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			desc:           "int lower than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			desc:           "int not lower than int64",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int32 lower than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: true,
		},
		{
			desc:           "int32 lower than int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int32(4), int32(5)},
			expectedResult: true,
		},
		{
			desc:           "int32 not lower than int32",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 lower than int",
			usrContext:     map[string]interface{}{"prop": int32(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not lower than int",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 lower than int64",
			usrContext:     map[string]interface{}{"prop": int32(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not lower than int64",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			desc:           "int64 lower than int64",
			usrContext:     map[string]interface{}{"prop": int64(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 lower than int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(5), int64(4)},
			expectedResult: true,
		},
		{
			desc:           "int64 not lower than int64",
			usrContext:     map[string]interface{}{"prop": int64(5)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			desc:           "int64 lower than int",
			usrContext:     map[string]interface{}{"prop": int64(0)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not lower than int",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int64 lower than int32",
			usrContext:     map[string]interface{}{"prop": int64(0)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not lower than int32",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			desc:           "float64 lower than float64",
			usrContext:     map[string]interface{}{"prop": float64(8.9)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 lower than float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(8.9)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 not lower than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.3)},
			property:       "prop",
			values:         []interface{}{float64(9.2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			res, err := operator.Lower(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestLowerOrEqual_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "int lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int(2)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int lower or equal than int from list",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int(3), int(4)},
			expectedResult: true,
		},
		{
			desc:           "int not lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: false,
		},
		{
			desc:           "int lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int(2)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		{
			desc:           "int not lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			desc:           "int lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			desc:           "int not lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int32 lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: true,
		},
		{
			desc:           "int32 lower or equal than int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int32(4), int32(5)},
			expectedResult: true,
		},
		{
			desc:           "int32 not lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int32(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int32(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			desc:           "int64 lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int64(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 lower or equal than int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(5), int64(4)},
			expectedResult: true,
		},
		{
			desc:           "int64 not lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int64(5)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			desc:           "int64 lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int64(0)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int64 lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int64(0)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			desc:           "float64 lower or equal than float64",
			usrContext:     map[string]interface{}{"prop": float64(8.9)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 lower or equal than float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(8.9)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 not lower or equal than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.3)},
			property:       "prop",
			values:         []interface{}{float64(9.2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			res, err := operator.LowerOrEqual(test.usrContext[test.property], test.values)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
