package flaggio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/flaggio/internal/flaggio"
)

func TestFlagCacheKey(t *testing.T) {
	tests := []struct {
		name        string
		parts       []string
		expectedKey string
	}{
		{
			name:        "uses no parts",
			parts:       []string{},
			expectedKey: "flaggio:flag",
		},
		{
			name:        "uses all parts",
			parts:       []string{"key", "my.flag.key"},
			expectedKey: "flaggio:flag:key:my.flag.key",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			key := flaggio.FlagCacheKey(tt.parts...)
			assert.Equal(t, tt.expectedKey, key)
		})
	}
}

func TestSegmentCacheKey(t *testing.T) {
	tests := []struct {
		name        string
		parts       []string
		expectedKey string
	}{
		{
			name:        "uses no parts",
			parts:       []string{},
			expectedKey: "flaggio:segment",
		},
		{
			name:        "uses all parts",
			parts:       []string{"key", "123"},
			expectedKey: "flaggio:segment:key:123",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			key := flaggio.SegmentCacheKey(tt.parts...)
			assert.Equal(t, tt.expectedKey, key)
		})
	}
}

func TestEvalCacheKey(t *testing.T) {
	tests := []struct {
		name        string
		parts       []string
		expectedKey string
	}{
		{
			name:        "uses no parts",
			parts:       []string{},
			expectedKey: "flaggio:eval",
		},
		{
			name:        "uses all parts",
			parts:       []string{"123", "456"},
			expectedKey: "flaggio:eval:123:456",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			key := flaggio.EvalCacheKey(tt.parts...)
			assert.Equal(t, tt.expectedKey, key)
		})
	}
}
