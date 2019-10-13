package repository

import (
	"context"

	"github.com/victorkohl/flaggio/internal/flaggio"
)

type Rule interface {
	CreateFlagRule(ctx context.Context, flagID string, input flaggio.NewFlagRule) (string, error)
	UpdateFlagRule(ctx context.Context, flagID, id string, input flaggio.UpdateFlagRule) error
	DeleteFlagRule(ctx context.Context, flagID, id string) error
}
