package operator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uw-labs/flaggio/internal/operator"
)

func TestInNetwork(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		usrContext     map[string]interface{}
		property       string
		values         []interface{}
		expectedResult bool
	}{
		{
			name:           "contains ipv4",
			usrContext:     map[string]interface{}{"$ip": "166.9.193.112"},
			property:       "$ip",
			values:         []interface{}{"166.9.193.0/8"},
			expectedResult: true,
		},
		{
			name:           "doesnt contain ipv4",
			usrContext:     map[string]interface{}{"$ip": "166.9.193.112"},
			property:       "$ip",
			values:         []interface{}{"166.9.194.0/24"},
			expectedResult: false,
		},
		{
			name:           "contains ipv6",
			usrContext:     map[string]interface{}{"$ip": "2001:0db8:0123:4567:89ab:cdef:1234:5678"},
			property:       "$ip",
			values:         []interface{}{"::/0"},
			expectedResult: true,
		},
		{
			name:           "doesnt contain ipv6",
			usrContext:     map[string]interface{}{"$ip": "2001:0db8:0123:4567:89ab:cdef:1234:5678"},
			property:       "$ip",
			values:         []interface{}{"2001:0db9:0:0:0:0:0:0/32"},
			expectedResult: false,
		},
		// ========================================================================
		{
			name:           "unknown config type",
			usrContext:     map[string]interface{}{"$ip": "abcdef"},
			property:       "$ip",
			values:         []interface{}{struct{}{}},
			expectedResult: false,
		},
		{
			name:           "non-string user type",
			usrContext:     map[string]interface{}{"$ip": nil},
			property:       "$ip",
			values:         []interface{}{"2001:0db9:0:0:0:0:0:0/32"},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := operator.InNetwork(tt.usrContext[tt.property], tt.values)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
