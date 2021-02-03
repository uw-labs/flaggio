package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/flaggio/internal/operator"
)

func TestOneOf(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "equals string",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"abc"},
			expectedResult: true,
		},
		{
			name:           "equals string from list",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "abc"},
			expectedResult: true,
		},
		{
			name:           "not equals string",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: false,
		},
		{
			name:           "nil not equals string",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{"abc"},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "equals []byte",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("abc")},
			expectedResult: true,
		},
		{
			name:           "equals []byte from list",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("other"), []byte("abc")},
			expectedResult: true,
		},
		{
			name:           "not equals []byte",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("cde")},
			expectedResult: false,
		},
		{
			name:           "nil not equals []byte",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{[]byte("abc")},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "equals bool",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{true},
			expectedResult: true,
		},
		{
			name:           "equals bool from list",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false, true},
			expectedResult: true,
		},
		{
			name:           "not equals bool",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false},
			expectedResult: false,
		},
		{
			name:           "nil not equals bool",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{true},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "int equals int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			name:           "int equals int from list",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0), int(1)},
			expectedResult: true,
		},
		{
			name:           "int not equals int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "nil not equals int",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "int equals int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			name:           "int not equals int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			name:           "int equals int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(1)},
			expectedResult: true,
		},
		{
			name:           "int not equals int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "int32 equals int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		{
			name:           "int32 equals int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(0), int32(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not equals int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			name:           "int32 equals int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not equals int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			name:           "int32 equals int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not equals int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "int64 equals int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			name:           "int64 equals int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(0), int64(3)},
			expectedResult: true,
		},
		{
			name:           "int64 not equals int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			name:           "int64 equals int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not equals int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "int64 equals int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not equals int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "uint equals uint",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(4)},
			expectedResult: true,
		},
		{
			name:           "uint equals uint from list",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(0), uint(4)},
			expectedResult: true,
		},
		{
			name:           "uint not equals uint",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: false,
		},
		{
			name:           "uint equals uint32",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(4)},
			expectedResult: true,
		},
		{
			name:           "uint not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: false,
		},
		{
			name:           "uint equals uint64",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(4)},
			expectedResult: true,
		},
		{
			name:           "uint not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "uint32 equals uint32",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: true,
		},
		{
			name:           "uint32 equals uint32 from list",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(0), uint32(5)},
			expectedResult: true,
		},
		{
			name:           "uint32 not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: false,
		},
		{
			name:           "uint32 equals uint",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: true,
		},
		{
			name:           "uint32 not equals uint",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: false,
		},
		{
			name:           "uint32 equals uint64",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: true,
		},
		{
			name:           "uint32 not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "uint64 equals uint64",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: true,
		},
		{
			name:           "uint64 equals uint64 from list",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(0), uint64(6)},
			expectedResult: true,
		},
		{
			name:           "uint64 not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(7)},
			expectedResult: false,
		},
		{
			name:           "uint64 equals uint32",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: true,
		},
		{
			name:           "uint64 not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(7)},
			expectedResult: false,
		},
		{
			name:           "uint64 equals uint",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: true,
		},
		{
			name:           "uint64 not equals uint",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(7)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "float64 equals float64",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 equals float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 not equals float64",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := operator.OneOf(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestNotOneOf(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "nil not equals string",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{"abc"},
			expectedResult: true,
		},
		{
			name:           "not equals string",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"abc"},
			expectedResult: false,
		},
		{
			name:           "not equals string from list",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"other", "strings", "abc"},
			expectedResult: false,
		},
		{
			name:           "equals string",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"cde"},
			expectedResult: true,
		},
		{
			name:           "equals string from list",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{"cde", "fgh"},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "nil not equals []byte",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{[]byte("abc")},
			expectedResult: true,
		},
		{
			name:           "not equals []byte",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("abc")},
			expectedResult: false,
		},
		{
			name:           "not equals []byte from list",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("other"), []byte("abc")},
			expectedResult: false,
		},
		{
			name:           "equals []byte",
			usrContext:     map[string]interface{}{"prop": []byte("abc")},
			property:       "prop",
			values:         []interface{}{[]byte("cde")},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "nil not equals bool",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{true},
			expectedResult: true,
		},
		{
			name:           "not equals bool",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{true},
			expectedResult: false,
		},
		{
			name:           "not equals bool from list",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false, true},
			expectedResult: false,
		},
		{
			name:           "equals bool",
			usrContext:     map[string]interface{}{"prop": true},
			property:       "prop",
			values:         []interface{}{false},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "nil is not equal int",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: true,
		},
		{
			name:           "int not equals int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: false,
		},
		{
			name:           "int not equals int from list",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0), int(1)},
			expectedResult: false,
		},
		{
			name:           "int equals int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			name:           "int not equals int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: false,
		},
		{
			name:           "int equals int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		{
			name:           "int not equals int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(1)},
			expectedResult: false,
		},
		{
			name:           "int equals int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "int32 not equals int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			name:           "int32 not equals int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(0), int32(2)},
			expectedResult: false,
		},
		{
			name:           "int32 equals int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: true,
		},
		{
			name:           "int32 equals int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "int32 equals int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: true,
		},
		{
			name:           "int32 not equals int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		{
			name:           "int32 equals int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "int64 not equals int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		{
			name:           "int64 not equals int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(0), int64(3)},
			expectedResult: false,
		},
		{
			name:           "int64 equals int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: true,
		},
		{
			name:           "int64 not equals int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: false,
		},
		{
			name:           "int64 equals int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			name:           "int64 not equals int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: false,
		},
		{
			name:           "int64 equals int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "uint not equals uint",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(4)},
			expectedResult: false,
		},
		{
			name:           "uint not equals uint from list",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(0), uint(4)},
			expectedResult: false,
		},
		{
			name:           "uint equals uint",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: true,
		},
		{
			name:           "uint not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(4)},
			expectedResult: false,
		},
		{
			name:           "uint equals uint32",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: true,
		},
		{
			name:           "uint not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(4)},
			expectedResult: false,
		},
		{
			name:           "uint equals uint64",
			usrContext:     map[string]interface{}{"prop": uint(4)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "uint32 not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(5)},
			expectedResult: false,
		},
		{
			name:           "uint32 not equals uint32 from list",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(0), uint32(5)},
			expectedResult: false,
		},
		{
			name:           "uint32 equals uint32",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: true,
		},
		{
			name:           "uint32 not equals uint",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(5)},
			expectedResult: false,
		},
		{
			name:           "uint32 equals uint",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: true,
		},
		{
			name:           "uint32 not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(5)},
			expectedResult: false,
		},
		{
			name:           "uint32 equals uint64",
			usrContext:     map[string]interface{}{"prop": uint32(5)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "uint64 not equals uint64",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(6)},
			expectedResult: false,
		},
		{
			name:           "uint64 not equals uint64 from list",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(0), uint64(6)},
			expectedResult: false,
		},
		{
			name:           "uint64 equals uint64",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint64(7)},
			expectedResult: true,
		},
		{
			name:           "uint64 not equals uint32",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(6)},
			expectedResult: false,
		},
		{
			name:           "uint64 equals uint32",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint32(7)},
			expectedResult: true,
		},
		{
			name:           "uint64 not equals uint",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(6)},
			expectedResult: false,
		},
		{
			name:           "uint64 equals uint",
			usrContext:     map[string]interface{}{"prop": uint64(6)},
			property:       "prop",
			values:         []interface{}{uint(7)},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "float64 not equals float64",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: false,
		},
		{
			name:           "float64 not equals float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: false,
		},
		{
			name:           "float64 not equals float64",
			usrContext:     map[string]interface{}{"prop": float64(9.1)},
			property:       "prop",
			values:         []interface{}{float64(9.2)},
			expectedResult: true,
		},
		// ========================================================================
		{
			name:           "unknown type",
			usrContext:     map[string]interface{}{"prop": "abc"},
			property:       "prop",
			values:         []interface{}{struct{}{}},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := operator.NotOneOf(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
