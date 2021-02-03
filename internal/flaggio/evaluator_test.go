package flaggio_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/uw-labs/flaggio/internal/flaggio"
	flaggio_mock "github.com/uw-labs/flaggio/internal/flaggio/mocks"
)

func TestEvaluate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		run  func(t *testing.T, mockCtrl *gomock.Controller)
	}{
		{
			name: "return correct answer when first evaluator has the answer",
			run: func(t *testing.T, mockCtrl *gomock.Controller) {
				userContext := map[string]interface{}{"name": "John"}
				eval := flaggio_mock.NewMockEvaluator(mockCtrl)
				res := flaggio.EvalResult{Answer: 1}

				eval.EXPECT().Evaluate(userContext).Times(1).Return(res, nil)

				result, err := flaggio.Evaluate(userContext, eval)
				assert.NoError(t, err)
				assert.Equal(t, 1, result.Answer)
				assert.Nil(t, result.Next)
			},
		},
		{
			name: "return correct answer when leaf evaluator in slice has the answer",
			run: func(t *testing.T, mockCtrl *gomock.Controller) {
				userContext := map[string]interface{}{"name": "John"}
				eval1 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval2 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval3 := flaggio_mock.NewMockEvaluator(mockCtrl)
				res1 := flaggio.EvalResult{Answer: 1, Next: []flaggio.Evaluator{eval2, eval3}}
				res2 := flaggio.EvalResult{}
				res3 := flaggio.EvalResult{Answer: "abc"}

				eval1.EXPECT().Evaluate(userContext).Times(1).Return(res1, nil)
				eval2.EXPECT().Evaluate(userContext).Times(1).Return(res2, nil)
				eval3.EXPECT().Evaluate(userContext).Times(1).Return(res3, nil)

				result, err := flaggio.Evaluate(userContext, eval1)
				assert.NoError(t, err)
				assert.Equal(t, "abc", result.Answer)
				assert.Nil(t, result.Next)
			},
		},
		{
			name: "return correct answer when leaf evaluator in chain has the answer",
			run: func(t *testing.T, mockCtrl *gomock.Controller) {
				userContext := map[string]interface{}{"name": "John"}
				eval1 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval2 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval3 := flaggio_mock.NewMockEvaluator(mockCtrl)
				res1 := flaggio.EvalResult{Answer: 1, Next: []flaggio.Evaluator{eval2}}
				res2 := flaggio.EvalResult{Next: []flaggio.Evaluator{eval3}}
				res3 := flaggio.EvalResult{Answer: true}

				eval1.EXPECT().Evaluate(userContext).Times(1).Return(res1, nil)
				eval2.EXPECT().Evaluate(userContext).Times(1).Return(res2, nil)
				eval3.EXPECT().Evaluate(userContext).Times(1).Return(res3, nil)

				result, err := flaggio.Evaluate(userContext, eval1)
				assert.NoError(t, err)
				assert.Equal(t, true, result.Answer)
				assert.Nil(t, result.Next)
			},
		},
		{
			name: "return correct answer when root evaluator has answer and chained evaluators in slice don't",
			run: func(t *testing.T, mockCtrl *gomock.Controller) {
				userContext := map[string]interface{}{"name": "John"}
				eval1 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval2 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval3 := flaggio_mock.NewMockEvaluator(mockCtrl)
				res1 := flaggio.EvalResult{Answer: "answer", Next: []flaggio.Evaluator{eval2, eval3}}
				res2 := flaggio.EvalResult{}
				res3 := flaggio.EvalResult{}

				eval1.EXPECT().Evaluate(userContext).Times(1).Return(res1, nil)
				eval2.EXPECT().Evaluate(userContext).Times(1).Return(res2, nil)
				eval3.EXPECT().Evaluate(userContext).Times(1).Return(res3, nil)

				result, err := flaggio.Evaluate(userContext, eval1)
				assert.NoError(t, err)
				assert.Equal(t, "answer", result.Answer)
				assert.Equal(t, []flaggio.Evaluator{eval2, eval3}, result.Next)
			},
		},
		{
			name: "return correct answer when deep chained evaluator has the answer",
			run: func(t *testing.T, mockCtrl *gomock.Controller) {
				userContext := map[string]interface{}{"name": "John"}
				eval1 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval2 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval3 := flaggio_mock.NewMockEvaluator(mockCtrl)
				res1 := flaggio.EvalResult{Next: []flaggio.Evaluator{eval2}}
				res2 := flaggio.EvalResult{Answer: 2, Next: []flaggio.Evaluator{eval3}}
				res3 := flaggio.EvalResult{}

				eval1.EXPECT().Evaluate(userContext).Times(1).Return(res1, nil)
				eval2.EXPECT().Evaluate(userContext).Times(1).Return(res2, nil)
				eval3.EXPECT().Evaluate(userContext).Times(1).Return(res3, nil)

				result, err := flaggio.Evaluate(userContext, eval1)
				assert.NoError(t, err)
				assert.Equal(t, 2, result.Answer)
				assert.Equal(t, []flaggio.Evaluator{eval3}, result.Next)
			},
		},
		{
			name: "return correct answer when deep chained evaluator in slice has the answer",
			run: func(t *testing.T, mockCtrl *gomock.Controller) {
				userContext := map[string]interface{}{"name": "John"}
				eval1 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval2 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval3 := flaggio_mock.NewMockEvaluator(mockCtrl)
				res1 := flaggio.EvalResult{Next: []flaggio.Evaluator{eval2, eval3}}
				res2 := flaggio.EvalResult{Answer: 5}
				res3 := flaggio.EvalResult{}

				eval1.EXPECT().Evaluate(userContext).Times(1).Return(res1, nil)
				eval2.EXPECT().Evaluate(userContext).Times(1).Return(res2, nil)
				eval3.EXPECT().Evaluate(userContext).Times(0).Return(res3, nil)

				result, err := flaggio.Evaluate(userContext, eval1)
				assert.NoError(t, err)
				assert.Equal(t, 5, result.Answer)
				assert.Nil(t, result.Next)
			},
		},
		{
			name: "return correct answer when root evaluator has answer and chained evaluators don't",
			run: func(t *testing.T, mockCtrl *gomock.Controller) {
				userContext := map[string]interface{}{"name": "John"}
				eval1 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval2 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval3 := flaggio_mock.NewMockEvaluator(mockCtrl)
				res1 := flaggio.EvalResult{Answer: 2.5, Next: []flaggio.Evaluator{eval2}}
				res2 := flaggio.EvalResult{Next: []flaggio.Evaluator{eval3}}
				res3 := flaggio.EvalResult{}

				eval1.EXPECT().Evaluate(userContext).Times(1).Return(res1, nil)
				eval2.EXPECT().Evaluate(userContext).Times(1).Return(res2, nil)
				eval3.EXPECT().Evaluate(userContext).Times(1).Return(res3, nil)

				result, err := flaggio.Evaluate(userContext, eval1)
				assert.NoError(t, err)
				assert.Equal(t, 2.5, result.Answer)
				assert.Equal(t, []flaggio.Evaluator{eval2}, result.Next)
			},
		},
		{
			name: "return nil answer when no evaluators have an answer",
			run: func(t *testing.T, mockCtrl *gomock.Controller) {
				userContext := map[string]interface{}{"name": "John"}
				eval1 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval2 := flaggio_mock.NewMockEvaluator(mockCtrl)
				eval3 := flaggio_mock.NewMockEvaluator(mockCtrl)
				res1 := flaggio.EvalResult{Next: []flaggio.Evaluator{eval2}}
				res2 := flaggio.EvalResult{Next: []flaggio.Evaluator{eval3}}
				res3 := flaggio.EvalResult{}

				eval1.EXPECT().Evaluate(userContext).Times(1).Return(res1, nil)
				eval2.EXPECT().Evaluate(userContext).Times(1).Return(res2, nil)
				eval3.EXPECT().Evaluate(userContext).Times(1).Return(res3, nil)

				result, err := flaggio.Evaluate(userContext, eval1)
				assert.NoError(t, err)
				assert.Nil(t, result.Answer)
				assert.Nil(t, result.Next)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			tt.run(t, mockCtrl)
		})
	}
}
