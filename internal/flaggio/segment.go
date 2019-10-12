package flaggio

import (
	"time"
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

func (s Segment) Validator(usrContext map[string]interface{}) bool {
	for _, rl := range s.Rules {
		if ConstraintList(rl.Constraints).Validate(usrContext) {
			return true
		}
	}
	return false
}
