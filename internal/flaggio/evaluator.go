package flaggio

import (
	"fmt"
)

type Evaluator interface {
	Evaluate(usr map[string]interface{}) (EvalResult, error)
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

type Trace struct {
	Type   string
	ID     string
	Answer interface{}
}

func (r EvalResult) Stack() (stack []Trace) {
	prev := &r
	for prev != nil {
		var id string
		ider, ok := prev.evaluator.(Identifier)
		if ok {
			id = ider.GetID()
		}
		trace := Trace{
			Type:   fmt.Sprintf("%T", prev.evaluator),
			ID:     id,
			Answer: prev.Answer,
		}
		stack = append(stack, trace)
		prev = prev.previous
	}
	return
}

func Evaluate(usr map[string]interface{}, root Evaluator) (EvalResult, error) {
	return evaluate(usr, []Evaluator{root})
}

func evaluate(usr map[string]interface{}, evaluators []Evaluator) (EvalResult, error) {
	var lastResult EvalResult
	var last *EvalResult
	for len(evaluators) > 0 {
		evltr := evaluators[0]
		res, err := evltr.Evaluate(usr)
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
