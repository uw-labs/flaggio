package flaggio

import (
	"time"

	"github.com/victorkt/flaggio/internal/operator"
)

var _ Identifier = (*Segment)(nil)
var _ operator.Validator = (*Segment)(nil)

// Segment represents a category of users that can be grouped based on a
// set of rules.
type Segment struct {
	ID          string
	Name        string
	Description *string
	Rules       []*SegmentRule
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

// GetID returns the segment ID.
func (s Segment) GetID() string {
	return s.ID
}

// Validate will check if any of the segment rules passes validation. If so,
// the validation is successful, otherwise it returns false.
func (s Segment) Validate(usrContext map[string]interface{}) (bool, error) {
	for _, rl := range s.Rules {
		ok, err := ConstraintList(rl.Constraints).Validate(usrContext)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
