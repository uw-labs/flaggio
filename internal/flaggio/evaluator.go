package flaggio

import (
	"fmt"
)

type Evaluator interface {
	Evaluate(usrContext map[string]interface{}) (EvalResult, error)
}

type Identifier interface {
	GetID() string
}

type EvalResult struct {
	Answer    interface{}
	Next      []Evaluator
	evaluator Evaluator
	previous  *EvalResult
}

func (r EvalResult) Stack() (stack []*StackTrace) {
	prev := &r
	for prev != nil {
		var id *string
		ider, ok := prev.evaluator.(Identifier)
		if ok {
			v := ider.GetID()
			id = &v
		}
		stack = append(stack, &StackTrace{
			Type:   fmt.Sprintf("%T", prev.evaluator),
			ID:     id,
			Answer: prev.Answer,
		})
		prev = prev.previous
	}
	return
}

func Evaluate(usrContext map[string]interface{}, root Evaluator) (EvalResult, error) {
	return evaluate(usrContext, []Evaluator{root})
}

func evaluate(usrContext map[string]interface{}, evaluators []Evaluator) (EvalResult, error) {
	var lastResult EvalResult
	var last *EvalResult
	for len(evaluators) > 0 {
		evltr := evaluators[0]
		res, err := evltr.Evaluate(usrContext)
		if err != nil {
			return EvalResult{}, err
		}
		res.evaluator = evltr
		res.previous = last
		last = &res
		if res.Answer != nil {
			if len(res.Next) == 0 {
				return res, nil
			}
			lastResult = res
		}
		if len(res.Next) > 0 {
			evaluators = append(res.Next, evaluators[1:]...)
		} else {
			evaluators = evaluators[1:]
		}
	}
	return lastResult, nil
}
