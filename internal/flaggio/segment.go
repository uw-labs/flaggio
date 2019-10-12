package flaggio

import (
	"time"

	"github.com/victorkohl/flaggio/internal/errors"
)

type Segment struct {
	ID          string
	Key         string
	Name        string
	Description *string
	Rules       []*SegmentRule
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func (r Segment) Evaluate(usr map[string]interface{}) (EvalResult, error) {
	// TODO: implement this
	return EvalResult{}, errors.ErrNotImplemented
}
