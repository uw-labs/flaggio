package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkohl/flaggio/internal/operator"
)

func TestOneOf_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usr            map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "equals string",
			usr:            map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"abc"},
			expectedResult: true,
		},
		{
			desc:           "equals string from list",
			usr:            map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "abc"},
			expectedResult: true,
		},
		{
			desc:           "not equals string",
			usr:            map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "equals []byte",
			usr:            map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("abc")},
			expectedResult: true,
		},
		{
			desc:           "equals []byte from list",
			usr:            map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("other"), []byte("abc")},
			expectedResult: true,
		},
		{
			desc:           "not equals []byte",
			usr:            map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("cde")},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "equals bool",
			usr:            map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{true},
			expectedResult: true,
		},
		{
			desc:           "equals bool from list",
			usr:            map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false, true},
			expectedResult: true,
		},
		{
			desc:           "not equals bool",
			usr:            map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int equals int",
			usr:            map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			desc:           "int equals int from list",
			usr:            map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0), int(1)},
			expectedResult: true,
		},
		{
			desc:           "int not equals int",
			usr:            map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int equals int32",
			usr:            map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int not equals int32",
			usr:            map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			desc:           "int equals int64",
			usr:            map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(1)},
			expectedResult: true,
		},
		{
			desc:           "int not equals int64",
			usr:            map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int32 equals int32",
			usr:            map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 equals int32 from list",
			usr:            map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(0), int32(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not equals int32",
			usr:            map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 equals int",
			usr:            map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not equals int",
			usr:            map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 equals int64",
			usr:            map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not equals int64",
			usr:            map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int64 equals int64",
			usr:            map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 equals int64 from list",
			usr:            map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(0), int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 not equals int64",
			usr:            map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			desc:           "int64 equals int",
			usr:            map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not equals int",
			usr:            map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int64 equals int32",
			usr:            map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not equals int32",
			usr:            map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "uint equals uint",
			usr:            map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(4)},
			expectedResult: true,
		},
		{
			desc:           "uint equals uint from list",
			usr:            map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(0), uint(4)},
			expectedResult: true,
		},
		{
			desc:           "uint not equals uint",
			usr:            map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: false,
		},
		{
			desc:           "uint equals uint32",
			usr:            map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(4)},
			expectedResult: true,
		},
		{
			desc:           "uint not equals uint32",
			usr:            map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: false,
		},
		{
			desc:           "uint equals uint64",
			usr:            map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(4)},
			expectedResult: true,
		},
		{
			desc:           "uint not equals uint64",
			usr:            map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "uint32 equals uint32",
			usr:            map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: true,
		},
		{
			desc:           "uint32 equals uint32 from list",
			usr:            map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(0), uint32(5)},
			expectedResult: true,
		},
		{
			desc:           "uint32 not equals uint32",
			usr:            map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: false,
		},
		{
			desc:           "uint32 equals uint",
			usr:            map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: true,
		},
		{
			desc:           "uint32 not equals uint",
			usr:            map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: false,
		},
		{
			desc:           "uint32 equals uint64",
			usr:            map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: true,
		},
		{
			desc:           "uint32 not equals uint64",
			usr:            map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "uint64 equals uint64",
			usr:            map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: true,
		},
		{
			desc:           "uint64 equals uint64 from list",
			usr:            map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(0), uint64(6)},
			expectedResult: true,
		},
		{
			desc:           "uint64 not equals uint64",
			usr:            map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(7)},
			expectedResult: false,
		},
		{
			desc:           "uint64 equals uint32",
			usr:            map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: true,
		},
		{
			desc:           "uint64 not equals uint32",
			usr:            map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(7)},
			expectedResult: false,
		},
		{
			desc:           "uint64 equals uint",
			usr:            map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: true,
		},
		{
			desc:           "uint64 not equals uint",
			usr:            map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(7)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "float32 equals float32",
			usr:            map[string]interface{}{"prop": float32(8.1)},
			property:       "prop",
			values:         []interface{}{float32(8.1)},
			expectedResult: true,
		},
		{
			desc:           "float32 equals float32 from list",
			usr:            map[string]interface{}{"prop": float32(8.1)},
			property:       "prop",
			values:         []interface{}{float32(8.0), float32(8.1)},
			expectedResult: true,
		},
		{
			desc:           "float32 not equals float32",
			usr:            map[string]interface{}{"prop": float32(8.1)},
			property:       "prop",
			values:         []interface{}{float32(8.2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "float64 equals float64",
			usr:            map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 equals float64 from list",
			usr:            map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 not equals float64",
			usr:            map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "unknown type",
			usr:            map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			op := operator.OneOf{}
			res := op.Operate(test.usr[test.property], test.values)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
