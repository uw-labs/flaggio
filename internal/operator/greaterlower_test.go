package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/flaggio/internal/operator"
)

func TestGreater(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "int greater than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: true,
		},
		{
			name:           "int greater than int from list",
			usrContext:     map[string]interface{}{"prop": int(2)},
			property:       "prop",
			values:         []interface{}{int(0), int(1)},
			expectedResult: true,
		},
		{
			name:           "int not greater than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "nil not greater than int",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: false,
		},
		{
			name:           "int greater than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(0)},
			expectedResult: true,
		},
		{
			name:           "int not greater than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			name:           "int greater than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(0)},
			expectedResult: true,
		},
		{
			name:           "int not greater than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "int32 greater than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			name:           "int32 greater than int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int32(0), int32(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not greater than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			name:           "int32 greater than int",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not greater than int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			name:           "int32 greater than int64",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not greater than int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			name:           "int64 greater than int64",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			name:           "int64 greater than int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int64(0), int64(3)},
			expectedResult: true,
		},
		{
			name:           "int64 not greater than int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			name:           "int64 greater than int",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not greater than int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "int64 greater than int32",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not greater than int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			name:           "float64 greater than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.2)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 greater than float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(9.2)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 not greater than float64",
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
			res, err := operator.Greater(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestGreaterOrEqual(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "int greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			name:           "int greater or equal than int from list",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0), int(1)},
			expectedResult: true,
		},
		{
			name:           "int not greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "nil not greater or equal than int",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: false,
		},
		{
			name:           "int greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(0)},
			expectedResult: true,
		},
		{
			name:           "int not greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			name:           "int greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(0)},
			expectedResult: true,
		},
		{
			name:           "int not greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "int32 greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			name:           "int32 greater or equal than int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int32(0), int32(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			name:           "int32 greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			name:           "int32 greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			name:           "int64 greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			name:           "int64 greater or equal than int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int64(0), int64(3)},
			expectedResult: true,
		},
		{
			name:           "int64 not greater or equal than int64",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			name:           "int64 greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not greater or equal than int",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "int64 greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int64(4)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not greater or equal than int32",
			usrContext:     map[string]interface{}{"prop": int64(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			name:           "float64 greater or equal than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.2)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 greater or equal than float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(9.2)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 not greater or equal than float64",
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
			res, err := operator.GreaterOrEqual(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestLower(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "int lower than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			name:           "int lower than int from list",
			usrContext:     map[string]interface{}{"prop": int(2)},
			property:       "prop",
			values:         []interface{}{int(3), int(4)},
			expectedResult: true,
		},
		{
			name:           "int not lower than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: false,
		},
		{
			name:           "nil not lower than int",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: false,
		},
		{
			name:           "int lower than int32",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		{
			name:           "int not lower than int32",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			name:           "int lower than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			name:           "int not lower than int64",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "int32 lower than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: true,
		},
		{
			name:           "int32 lower than int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int32(4), int32(5)},
			expectedResult: true,
		},
		{
			name:           "int32 not lower than int32",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			name:           "int32 lower than int",
			usrContext:     map[string]interface{}{"prop": int32(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not lower than int",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			name:           "int32 lower than int64",
			usrContext:     map[string]interface{}{"prop": int32(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not lower than int64",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			name:           "int64 lower than int64",
			usrContext:     map[string]interface{}{"prop": int64(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			name:           "int64 lower than int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(5), int64(4)},
			expectedResult: true,
		},
		{
			name:           "int64 not lower than int64",
			usrContext:     map[string]interface{}{"prop": int64(5)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			name:           "int64 lower than int",
			usrContext:     map[string]interface{}{"prop": int64(0)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not lower than int",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "int64 lower than int32",
			usrContext:     map[string]interface{}{"prop": int64(0)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not lower than int32",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			name:           "float64 lower than float64",
			usrContext:     map[string]interface{}{"prop": float64(8.9)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 lower than float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(8.9)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 not lower than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.3)},
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
			res, err := operator.Lower(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestLowerOrEqual(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "int lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int(2)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			name:           "int lower or equal than int from list",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int(3), int(4)},
			expectedResult: true,
		},
		{
			name:           "int not lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: false,
		},
		{
			name:           "nil not lower or equal than int",
			usrContext:     map[string]interface{}{"prop": nil},
			property:       "prop",
			values:         []interface{}{int(0)},
			expectedResult: false,
		},
		{
			name:           "int lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int(2)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: true,
		},
		{
			name:           "int not lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		{
			name:           "int lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			name:           "int not lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int(3)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "int32 lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int32(2)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: true,
		},
		{
			name:           "int32 lower or equal than int32 from list",
			usrContext:     map[string]interface{}{"prop": int32(3)},
			property:       "prop",
			values:         []interface{}{int32(4), int32(5)},
			expectedResult: true,
		},
		{
			name:           "int32 not lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int32(3)},
			expectedResult: false,
		},
		{
			name:           "int32 lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int32(1)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int(3)},
			expectedResult: false,
		},
		{
			name:           "int32 lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int32(1)},
			property:       "prop",
			values:         []interface{}{int64(2)},
			expectedResult: true,
		},
		{
			name:           "int32 not lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int32(4)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			name:           "int64 lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int64(2)},
			property:       "prop",
			values:         []interface{}{int64(3)},
			expectedResult: true,
		},
		{
			name:           "int64 lower or equal than int64 from list",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int64(5), int64(4)},
			expectedResult: true,
		},
		{
			name:           "int64 not lower or equal than int64",
			usrContext:     map[string]interface{}{"prop": int64(5)},
			property:       "prop",
			values:         []interface{}{int64(4)},
			expectedResult: false,
		},
		{
			name:           "int64 lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int64(0)},
			property:       "prop",
			values:         []interface{}{int(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not lower or equal than int",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int(2)},
			expectedResult: false,
		},
		{
			name:           "int64 lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int64(0)},
			property:       "prop",
			values:         []interface{}{int32(1)},
			expectedResult: true,
		},
		{
			name:           "int64 not lower or equal than int32",
			usrContext:     map[string]interface{}{"prop": int64(3)},
			property:       "prop",
			values:         []interface{}{int32(2)},
			expectedResult: false,
		},
		// // ========================================================================
		{
			name:           "float64 lower or equal than float64",
			usrContext:     map[string]interface{}{"prop": float64(8.9)},
			property:       "prop",
			values:         []interface{}{float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 lower or equal than float64 from list",
			usrContext:     map[string]interface{}{"prop": float64(8.9)},
			property:       "prop",
			values:         []interface{}{float64(9.0), float64(9.1)},
			expectedResult: true,
		},
		{
			name:           "float64 not lower or equal than float64",
			usrContext:     map[string]interface{}{"prop": float64(9.3)},
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
			res, err := operator.LowerOrEqual(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
