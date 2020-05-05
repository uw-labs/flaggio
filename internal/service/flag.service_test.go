package service_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	repository_mock "github.com/victorkt/flaggio/internal/repository/mocks"
	"github.com/victorkt/flaggio/internal/service"
)

var (
	ctxInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
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
		flagResult         *flaggio.Flag
		evaluationResult   *flaggio.Evaluation
		evaluationRequest  *service.EvaluationRequest
		expectedEvaluation *service.EvaluationResponse
		shouldReplaceEval  bool
	}{
		{
			name:       "return correct evaluation without debug option",
			flagKey:    "a",
			flagResult: flags[0],
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user1",
				UserContext: flaggio.UserContext{"name": "John"},
			},
			expectedEvaluation: &service.EvaluationResponse{
				Evaluation: &flaggio.Evaluation{FlagID: "1", FlagKey: "a", Value: 20,
					RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
			},
			shouldReplaceEval: true,
		},
		{
			name:       "return correct evaluation with debug option",
			flagKey:    "b",
			flagResult: flags[1],
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user2",
				UserContext: flaggio.UserContext{"name": "John"},
				Debug:       boolPtr(true),
			},
			expectedEvaluation: &service.EvaluationResponse{
				Evaluation: &flaggio.Evaluation{FlagID: "2", FlagKey: "b", Value: 10,
					RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55", StackTrace: []*flaggio.StackTrace{
						{Type: "*Flag", ID: stringPtr("2"), Answer: 10},
					}},
				UserContext: &flaggio.UserContext{"name": "John"},
			},
			shouldReplaceEval: false,
		},
		{
			name:       "return previous evaluation",
			flagKey:    "b",
			flagResult: flags[1],
			evaluationResult: &flaggio.Evaluation{FlagID: "1", FlagKey: "a", Value: 10,
				RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user2",
				UserContext: flaggio.UserContext{"name": "John"},
			},
			expectedEvaluation: &service.EvaluationResponse{
				Evaluation: &flaggio.Evaluation{FlagID: "1", FlagKey: "a", Value: 10,
					RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
			},
			shouldReplaceEval: false,
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
			evalRepo := repository_mock.NewMockEvaluation(mockCtrl)
			userRepo := repository_mock.NewMockUser(mockCtrl)
			flagService := service.NewFlagService(flagRepo, segmentRepo, evalRepo, userRepo)
			segmentResults := make([]*flaggio.Segment, 0)
			hash, err := tt.evaluationRequest.Hash()
			assert.NoError(t, err)

			flagRepo.EXPECT().
				FindByKey(gomock.AssignableToTypeOf(ctxInterface), tt.flagKey).
				Times(1).Return(tt.flagResult, nil)
			if tt.evaluationResult == nil {
				segmentRepo.EXPECT().
					FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, nil).
					Times(1).Return(segmentResults, nil)
			}
			evalRepo.EXPECT().
				FindByReqHashAndFlagKey(gomock.AssignableToTypeOf(ctxInterface), hash, tt.flagResult.Key).
				Times(1).Return(tt.evaluationResult, nil)
			if tt.shouldReplaceEval {
				userRepo.EXPECT().
					Replace(gomock.AssignableToTypeOf(ctxInterface), tt.evaluationRequest.UserID, tt.evaluationRequest.UserContext).
					Times(1).Return(nil)
				evalRepo.EXPECT().
					ReplaceOne(gomock.AssignableToTypeOf(ctxInterface), tt.evaluationRequest.UserID, tt.expectedEvaluation.Evaluation).
					Times(1).Return(nil)
			}

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
		evaluationResults  *flaggio.EvaluationResults
		expectedEvaluation *service.EvaluationsResponse
		outdatedEvals      flaggio.EvaluationList
		shouldReplaceEval  bool
	}{
		{
			name: "return correct evaluation without debug option",
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user1",
				UserContext: flaggio.UserContext{"name": "John"},
			},
			evaluationResults: &flaggio.EvaluationResults{Evaluations: []*flaggio.Evaluation{}},
			expectedEvaluation: &service.EvaluationsResponse{
				Evaluations: flaggio.EvaluationList{
					{FlagID: "1", FlagKey: "a", Value: 20, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
					{FlagID: "2", FlagKey: "b", Value: 10, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
				},
			},
			outdatedEvals: flaggio.EvaluationList{
				{FlagID: "1", FlagKey: "a", Value: 20, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
				{FlagID: "2", FlagKey: "b", Value: 10, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
			},
			shouldReplaceEval: true,
		},
		{
			name: "return correct evaluation with debug option",
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user2",
				UserContext: flaggio.UserContext{"name": "John"},
				Debug:       boolPtr(true),
			},
			evaluationResults: &flaggio.EvaluationResults{Evaluations: []*flaggio.Evaluation{}},
			expectedEvaluation: &service.EvaluationsResponse{
				Evaluations: flaggio.EvaluationList{
					{FlagID: "1", FlagKey: "a", Value: 20, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
					{FlagID: "2", FlagKey: "b", Value: 10, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
				},
				UserContext: &flaggio.UserContext{"name": "John"},
			},
			shouldReplaceEval: false,
		},
		{
			name: "return previous evaluation",
			evaluationRequest: &service.EvaluationRequest{
				UserID:      "user2",
				UserContext: flaggio.UserContext{"name": "John"},
			},
			evaluationResults: &flaggio.EvaluationResults{Evaluations: []*flaggio.Evaluation{
				{FlagID: "1", FlagKey: "a", Value: 10, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
			}},
			expectedEvaluation: &service.EvaluationsResponse{
				Evaluations: flaggio.EvaluationList{
					{FlagID: "1", FlagKey: "a", Value: 10, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
					{FlagID: "2", FlagKey: "b", Value: 10, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
				},
			},
			outdatedEvals: flaggio.EvaluationList{
				{FlagID: "2", FlagKey: "b", Value: 10, RequestHash: "5e83501f42ab66e04cd03a53d55399ffa7387a55"},
			},
			shouldReplaceEval: true,
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
			evalRepo := repository_mock.NewMockEvaluation(mockCtrl)
			userRepo := repository_mock.NewMockUser(mockCtrl)
			flagService := service.NewFlagService(flagRepo, segmentRepo, evalRepo, userRepo)
			flagResults := &flaggio.FlagResults{Flags: flags, Total: len(flags)}
			segmentResults := make([]*flaggio.Segment, 0)
			hash, err := tt.evaluationRequest.Hash()
			assert.NoError(t, err)

			flagRepo.EXPECT().
				FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, nil, nil).
				Times(1).Return(flagResults, nil)
			segmentRepo.EXPECT().
				FindAll(gomock.AssignableToTypeOf(ctxInterface), nil, nil).
				Times(1).Return(segmentResults, nil)
			evalRepo.EXPECT().
				FindAllByReqHash(gomock.AssignableToTypeOf(ctxInterface), hash).
				Times(1).Return(tt.evaluationResults.Evaluations, nil)
			if tt.shouldReplaceEval {
				userRepo.EXPECT().
					Replace(gomock.AssignableToTypeOf(ctxInterface), tt.evaluationRequest.UserID, tt.evaluationRequest.UserContext).
					Times(1).Return(nil)
				evalRepo.EXPECT().
					ReplaceAll(gomock.AssignableToTypeOf(ctxInterface), tt.evaluationRequest.UserID, hash, tt.outdatedEvals).
					Times(1).Return(nil)
			}

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
