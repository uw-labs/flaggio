package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/operator"
)

func TestOneOf_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "equals string",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"abc"},
			expectedResult: true,
		},
		{
			desc:           "equals string from list",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "abc"},
			expectedResult: true,
		},
		{
			desc:           "not equals string",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "equals []byte",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("abc")},
			expectedResult: true,
		},
		{
			desc:           "equals []byte from list",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("other"), []byte("abc")},
			expectedResult: true,
		},
		{
			desc:           "not equals []byte",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("cde")},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "equals bool",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{true},
			expectedResult: true,
		},
		{
			desc:           "equals bool from list",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false, true},
			expectedResult: true,
		},
		{
			desc:           "not equals bool",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int equals int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			desc:           "int equals int from list",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0), int(1)},
			expectedResult: true,
		},
		{
			desc:           "int not equals int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int equals int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int not equals int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			desc:           "int equals int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(1)},
			expectedResult: true,
		},
		{
			desc:           "int not equals int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int32 equals int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 equals int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(0), int32(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not equals int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 equals int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not equals int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			desc:           "int32 equals int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			desc:           "int32 not equals int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "int64 equals int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 equals int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(0), int64(3)},
			expectedResult: true,
		},
		{
			desc:           "int64 not equals int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			desc:           "int64 equals int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not equals int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int64 equals int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			desc:           "int64 not equals int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "uint equals uint",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(4)},
			expectedResult: true,
		},
		{
			desc:           "uint equals uint from list",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(0), uint(4)},
			expectedResult: true,
		},
		{
			desc:           "uint not equals uint",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: false,
		},
		{
			desc:           "uint equals uint32",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(4)},
			expectedResult: true,
		},
		{
			desc:           "uint not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: false,
		},
		{
			desc:           "uint equals uint64",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(4)},
			expectedResult: true,
		},
		{
			desc:           "uint not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "uint32 equals uint32",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: true,
		},
		{
			desc:           "uint32 equals uint32 from list",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(0), uint32(5)},
			expectedResult: true,
		},
		{
			desc:           "uint32 not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: false,
		},
		{
			desc:           "uint32 equals uint",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: true,
		},
		{
			desc:           "uint32 not equals uint",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: false,
		},
		{
			desc:           "uint32 equals uint64",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: true,
		},
		{
			desc:           "uint32 not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "uint64 equals uint64",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: true,
		},
		{
			desc:           "uint64 equals uint64 from list",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(0), uint64(6)},
			expectedResult: true,
		},
		{
			desc:           "uint64 not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(7)},
			expectedResult: false,
		},
		{
			desc:           "uint64 equals uint32",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: true,
		},
		{
			desc:           "uint64 not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(7)},
			expectedResult: false,
		},
		{
			desc:           "uint64 equals uint",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: true,
		},
		{
			desc:           "uint64 not equals uint",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(7)},
			expectedResult: false,
		},
		// ========================================================================
		{
			desc:           "float64 equals float64",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 equals float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			desc:           "float64 not equals float64",
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
			op := operator.OneOf{}
			res := op.Operate(test.usrContext[test.property], test.values)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func TestNotOneOf_Operate(t *testing.T) {
	tt := []struct {
		desc           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			desc:           "not equals string",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"abc"},
			expectedResult: false,
		},
		{
			desc:           "not equals string from list",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "abc"},
			expectedResult: false,
		},
		{
			desc:           "equals string",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: true,
		},
		{
			desc:           "equals string from list",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"cde", "fgh"},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "not equals []byte",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("abc")},
			expectedResult: false,
		},
		{
			desc:           "not equals []byte from list",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("other"), []byte("abc")},
			expectedResult: false,
		},
		{
			desc:           "equals []byte",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("cde")},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "not equals bool",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{true},
			expectedResult: false,
		},
		{
			desc:           "not equals bool from list",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false, true},
			expectedResult: false,
		},
		{
			desc:           "equals bool",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "int not equals int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: false,
		},
		{
			desc:           "int not equals int from list",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0), int(1)},
			expectedResult: false,
		},
		{
			desc:           "int equals int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int not equals int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: false,
		},
		{
			desc:           "int equals int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		{
			desc:           "int not equals int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(1)},
			expectedResult: false,
		},
		{
			desc:           "int equals int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "int32 not equals int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			desc:           "int32 not equals int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(0), int32(2)},
			expectedResult: false,
		},
		{
			desc:           "int32 equals int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: true,
		},
		{
			desc:           "int32 equals int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			desc:           "int32 equals int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: true,
		},
		{
			desc:           "int32 not equals int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		{
			desc:           "int32 equals int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "int64 not equals int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		{
			desc:           "int64 not equals int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(0), int64(3)},
			expectedResult: false,
		},
		{
			desc:           "int64 equals int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: true,
		},
		{
			desc:           "int64 not equals int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: false,
		},
		{
			desc:           "int64 equals int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			desc:           "int64 not equals int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: false,
		},
		{
			desc:           "int64 equals int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "uint not equals uint",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(4)},
			expectedResult: false,
		},
		{
			desc:           "uint not equals uint from list",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(0), uint(4)},
			expectedResult: false,
		},
		{
			desc:           "uint equals uint",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: true,
		},
		{
			desc:           "uint not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(4)},
			expectedResult: false,
		},
		{
			desc:           "uint equals uint32",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: true,
		},
		{
			desc:           "uint not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(4)},
			expectedResult: false,
		},
		{
			desc:           "uint equals uint64",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "uint32 not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: false,
		},
		{
			desc:           "uint32 not equals uint32 from list",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(0), uint32(5)},
			expectedResult: false,
		},
		{
			desc:           "uint32 equals uint32",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: true,
		},
		{
			desc:           "uint32 not equals uint",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: false,
		},
		{
			desc:           "uint32 equals uint",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: true,
		},
		{
			desc:           "uint32 not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: false,
		},
		{
			desc:           "uint32 equals uint64",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "uint64 not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: false,
		},
		{
			desc:           "uint64 not equals uint64 from list",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(0), uint64(6)},
			expectedResult: false,
		},
		{
			desc:           "uint64 equals uint64",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(7)},
			expectedResult: true,
		},
		{
			desc:           "uint64 not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: false,
		},
		{
			desc:           "uint64 equals uint32",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(7)},
			expectedResult: true,
		},
		{
			desc:           "uint64 not equals uint",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: false,
		},
		{
			desc:           "uint64 equals uint",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(7)},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "float64 not equals float64",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: false,
		},
		{
			desc:           "float64 not equals float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: false,
		},
		{
			desc:           "float64 not equals float64",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.2)},
			expectedResult: true,
		},
		// ========================================================================
		{
			desc:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: true,
		},
	}

	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			op := operator.NotOneOf{}
			res := op.Operate(test.usrContext[test.property], test.values)
			assert.Equal(t, test.expectedResult, res)
		})
	}
}
