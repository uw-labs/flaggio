package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	repository_mock "github.com/victorkt/flaggio/internal/repository/mocks"
	"github.com/victorkt/flaggio/internal/server/api/service"
)

func TestFlagService_Evaluate(t *testing.T) {
	t.Parallel()
	variants := []*flaggio.Variant{
		{ID: "1", Value: 10},
		{ID: "2", Value: 20},
	}
	flags := []*flaggio.Flag{
		{ID: "1", Key: "a", Enabled: false, Variants: variants, DefaultVariantWhenOn: variants[0], DefaultVariantWhenOff: variants[1]},
		{ID: "2", Key: "b", Enabled: true, Variants: variants, DefaultVariantWhenOn: variants[0], DefaultVariantWhenOff: variants[1]},
	}
	tests := []struct {
		name               string
		flagKey            string
		evaluationRequest  *service.EvaluationRequest
		expectedEvaluation *service.EvaluationResponse
		run                func(t *testing.T, mockCtrl *gomock.Controller)
	}{
		{
			name:    "return correct evaluation without debug option",
			flagKey: "a",
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user1",
				UserContext: flaggio.UserContext{"name": "John"},
			},
			expectedEvaluation: &service.EvaluationResponse{
				Evaluation: &flaggio.Evaluation{FlagKey: "a", Value: 20},
			},
		},
		{
			name:    "return correct evaluation with debug option",
			flagKey: "b",
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user2",
				UserContext: flaggio.UserContext{"name": "John"},
				Debug:       boolPtr(true),
			},
			expectedEvaluation: &service.EvaluationResponse{
				Evaluation: &flaggio.Evaluation{FlagKey: "b", Value: 20, StackTrace: []*flaggio.StackTrace{
					{Type: "*Flag", ID: stringPtr("1"), Answer: 20},
				}},
				UserContext: &flaggio.UserContext{"name": "John"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ctx := context.Background()
			flagRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentRepo := repository_mock.NewMockSegment(mockCtrl)
			flagService := service.NewFlagService(flagRepo, segmentRepo)
			flagResults := flags[0]
			segmentesults := make([]*flaggio.Segment, 0)

			flagRepo.EXPECT().
				FindByKey(ctx, tt.flagKey).
				Times(1).Return(flagResults, nil)
			segmentRepo.EXPECT().
				FindAll(ctx, nil, nil).
				Times(1).Return(segmentesults, nil)

			result, err := flagService.Evaluate(ctx, tt.flagKey, tt.evaluationRequest)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedEvaluation, result)
		})
	}
}

func TestFlagService_EvaluateAll(t *testing.T) {
	t.Parallel()
	variants := []*flaggio.Variant{
		{ID: "1", Value: 10},
		{ID: "2", Value: 20},
	}
	flags := []*flaggio.Flag{
		{ID: "1", Key: "a", Enabled: false, Variants: variants, DefaultVariantWhenOn: variants[0], DefaultVariantWhenOff: variants[1]},
		{ID: "2", Key: "b", Enabled: true, Variants: variants, DefaultVariantWhenOn: variants[0], DefaultVariantWhenOff: variants[1]},
	}
	tests := []struct {
		name               string
		evaluationRequest  *service.EvaluationRequest
		expectedEvaluation *service.EvaluationsResponse
		run                func(t *testing.T, mockCtrl *gomock.Controller)
	}{
		{
			name: "return correct evaluation without debug option",
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user1",
				UserContext: flaggio.UserContext{"name": "John"},
			},
			expectedEvaluation: &service.EvaluationsResponse{
				Evaluations: flaggio.EvaluationList{
					{FlagKey: "a", Value: 20},
					{FlagKey: "b", Value: 10},
				},
			},
		},
		{
			name: "return correct evaluation with debug option",
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user2",
				UserContext: flaggio.UserContext{"name": "John"},
				Debug:       boolPtr(true),
			},
			expectedEvaluation: &service.EvaluationsResponse{
				Evaluations: flaggio.EvaluationList{
					{FlagKey: "a", Value: 20},
					{FlagKey: "b", Value: 10},
				},
				UserContext: &flaggio.UserContext{"name": "John"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ctx := context.Background()
			flagRepo := repository_mock.NewMockFlag(mockCtrl)
			segmentRepo := repository_mock.NewMockSegment(mockCtrl)
			flagService := service.NewFlagService(flagRepo, segmentRepo)
			flagResults := &flaggio.FlagResults{Flags: flags, Total: len(flags)}
			segmentesults := make([]*flaggio.Segment, 0)

			flagRepo.EXPECT().
				FindAll(ctx, nil, nil, nil).
				Times(1).Return(flagResults, nil)
			segmentRepo.EXPECT().
				FindAll(ctx, nil, nil).
				Times(1).Return(segmentesults, nil)

			result, err := flagService.EvaluateAll(ctx, tt.evaluationRequest)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedEvaluation, result)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
