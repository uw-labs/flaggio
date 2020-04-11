package redis_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/service"
	service_mock "github.com/victorkt/flaggio/internal/service/mocks"
	redis_svc "github.com/victorkt/flaggio/internal/service/redis"
)

var (
	evalsList = flaggio.EvaluationList{
		{FlagKey: "f1", Value: "abc"},
		{FlagKey: "f2", Value: int64(1)},
	}
	ctxInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func TestFlagService_Evaluate(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(t *testing.T, repo *service_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying service on cache miss",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalRes := &service.EvaluationResponse{Evaluation: evalsList[0], UserContext: nil}
				flagRedisSvc := redis_svc.NewFlagService(redisClient, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
				}
				flagSvc.EXPECT().Evaluate(gomock.AssignableToTypeOf(ctxInterface), "f1", evalReq).
					Times(1).Return(evalRes, nil)

				res, err := flagRedisSvc.Evaluate(ctx, "f1", evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalRes, res)
			},
		},
		{
			name: "doesnt call underlying repository on cache hit",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalRes := &service.EvaluationResponse{Evaluation: evalsList[0], UserContext: nil}
				flagRedisSvc := redis_svc.NewFlagService(redisClient, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
				}
				flagSvc.EXPECT().Evaluate(gomock.AssignableToTypeOf(ctxInterface), "f1", evalReq).
					Times(0)

				res, err := flagRedisSvc.Evaluate(ctx, "f1", evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalRes, res)
			},
		},
		{
			name: "always calls underlying service on debug requests",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalRes := &service.EvaluationResponse{
					Evaluation:  evalsList[0],
					UserContext: &flaggio.UserContext{"test": "abc"}}
				flagRedisSvc := redis_svc.NewFlagService(redisClient, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
					Debug:       boolPtr(true),
				}
				flagSvc.EXPECT().Evaluate(gomock.AssignableToTypeOf(ctxInterface), "f1", evalReq).
					Times(1).Return(evalRes, nil)

				res, err := flagRedisSvc.Evaluate(ctx, "f1", evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalRes, res)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			flagSvc := service_mock.NewMockFlag(mockCtrl)

			tt.run(t, flagSvc)
		})
	}
}

func TestFlagService_EvaluateAll(t *testing.T) {
	// flush cache first
	if err := redisClient.FlushAll().Err(); err != nil {
		t.Fatalf("failed to flush cache: %s", err)
	}

	tests := []struct {
		name string
		run  func(t *testing.T, repo *service_mock.MockFlag)
	}{
		// these tests are meant to be run in order
		{
			name: "calls underlying service on cache miss",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalsRes := &service.EvaluationsResponse{Evaluations: evalsList, UserContext: nil}
				flagRedisSvc := redis_svc.NewFlagService(redisClient, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
				}
				flagSvc.EXPECT().EvaluateAll(gomock.AssignableToTypeOf(ctxInterface), evalReq).
					Times(1).Return(evalsRes, nil)

				res, err := flagRedisSvc.EvaluateAll(ctx, evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalsRes, res)
			},
		},
		{
			name: "doesnt call underlying repository on cache hit",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalsRes := &service.EvaluationsResponse{Evaluations: evalsList, UserContext: nil}
				flagRedisSvc := redis_svc.NewFlagService(redisClient, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
				}
				flagSvc.EXPECT().EvaluateAll(gomock.AssignableToTypeOf(ctxInterface), evalReq).
					Times(0)

				res, err := flagRedisSvc.EvaluateAll(ctx, evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalsRes, res)
			},
		},
		{
			name: "always calls underlying service on debug requests",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalsRes := &service.EvaluationsResponse{
					Evaluations: evalsList,
					UserContext: &flaggio.UserContext{"test": "abc"}}
				flagRedisSvc := redis_svc.NewFlagService(redisClient, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
					Debug:       boolPtr(true),
				}
				flagSvc.EXPECT().EvaluateAll(gomock.AssignableToTypeOf(ctxInterface), evalReq).
					Times(1).Return(evalsRes, nil)

				res, err := flagRedisSvc.EvaluateAll(ctx, evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalsRes, res)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			flagSvc := service_mock.NewMockFlag(mockCtrl)

			tt.run(t, flagSvc)
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
