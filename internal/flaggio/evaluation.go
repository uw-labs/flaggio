package flaggio

import (
	"net/http"
)

type Evaluation struct {
	FlagKey    string        `json:"flagKey"`
	Value      interface{}   `json:"value,omitempty"`
	Error      string        `json:"error,omitempty"`
	StackTrace []*StackTrace `json:"stackTrace,omitempty"`
}

func (e *Evaluation) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type StackTrace struct {
	Type   string      `json:"type"`
	ID     *string     `json:"id"`
	Answer interface{} `json:"answer"`
}

type EvaluationList []*Evaluation

func (l EvaluationList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
