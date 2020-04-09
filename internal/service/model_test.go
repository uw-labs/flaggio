package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/service"
)

func TestEvaluationRequest_Hash(t *testing.T) {
	tests := []struct {
		name         string
		req          service.EvaluationRequest
		expectedHash string
	}{
		{
			name: "returns hash",
			req: service.EvaluationRequest{
				UserID: "123",
				UserContext: flaggio.UserContext{
					"abc": 123,
					"cde": "456",
					"efg": true,
					"ghi": nil,
				},
				Debug: nil,
			},
			expectedHash: "78c77cefc3d6a062e7c29140c4aef97be2d8e0c4",
		},
		{
			name: "returns same hash regardless of order",
			req: service.EvaluationRequest{
				UserID: "123",
				UserContext: flaggio.UserContext{
					"cde": "456",
					"abc": 123,
					"ghi": nil,
					"efg": true,
				},
				Debug: nil,
			},
			expectedHash: "78c77cefc3d6a062e7c29140c4aef97be2d8e0c4",
		},
		{
			name: "doesnt care about debug property",
			req: service.EvaluationRequest{
				UserID: "123",
				UserContext: flaggio.UserContext{
					"cde": "456",
					"abc": 123,
					"ghi": nil,
					"efg": true,
				},
				Debug: boolPtr(true),
			},
			expectedHash: "78c77cefc3d6a062e7c29140c4aef97be2d8e0c4",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			hash, err := tt.req.Hash()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedHash, hash)
		})
	}
}
