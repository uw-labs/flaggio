package flaggio

import (
	"net/http"
	"time"
)

// Evaluation is the final result of a flag evaluation. It holds the
// returned value associated with the key for the given user.
// If an error occurred, value will be nil and the error property will
// contain the error message.
// Optionally, a stack trace of the evaluation process can be attached
// to the object.
type Evaluation struct {
	ID          string        `json:"-"`
	FlagID      string        `json:"-"`
	FlagVersion int           `json:"-"`
	RequestHash string        `json:"-"`
	CreatedAt   time.Time     `json:"-"`
	FlagKey     string        `json:"flagKey"`
	Value       interface{}   `json:"value,omitempty"`
	Error       string        `json:"error,omitempty"`
	StackTrace  []*StackTrace `json:"stackTrace,omitempty"`
}

// StackTrace contains detailed information about the evaluation process.
// Type is the type of the model object that evaluated the user context
// ID holds the ID of the same object, if any. Answer is the evaluation
// answer, if any.
type StackTrace struct {
	Type   string      `json:"type"`
	ID     *string     `json:"id"`
	Answer interface{} `json:"answer"`
}

// EvaluationList is a slice of *Evaluation.
type EvaluationList []*Evaluation

// Render can enrich the EvaluationList before being returned to the
// user. Currently it does nothing, but is needed to satisfy the
// chi.Renderer interface.
func (l EvaluationList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
