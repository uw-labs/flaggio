package flaggio

import (
	"time"

	"github.com/victorkt/flaggio/internal/operator"
)

var _ Identifier = (*Segment)(nil)
var _ operator.Validator = (*Segment)(nil)

type Segment struct {
	ID          string
	Name        string
	Description *string
	Rules       []*SegmentRule
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func (s Segment) GetID() string {
	return s.ID
}

func (s Segment) Validate(usrContext map[string]interface{}) bool {
	for _, rl := range s.Rules {
		if ConstraintList(rl.Constraints).Validate(usrContext) {
			return true
		}
	}
	return false
}
