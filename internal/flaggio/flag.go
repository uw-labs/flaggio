package flaggio

import (
	"time"
)

type Flag struct {
	ID          string
	Key         string
	Name        string
	Description *string
	Enabled     bool
	Version     uint64
	Variants    []*Variant
	Rules       []*FlagRule
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
