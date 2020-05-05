package repository

//go:generate mockgen -destination=./mocks/user_mock.go -package=repository_mock github.com/victorkt/flaggio/internal/repository User

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

// User represents a set of operations available to list and manage users.
type User interface {
	// FindAll returns a list of users, based on an optional offset and limit.
	FindAll(ctx context.Context, search *string, offset, limit *int64) (*flaggio.UserResults, error)
	// FindByID returns a user by its id.
	FindByID(ctx context.Context, id string) (*flaggio.User, error)
	// Replace creates or updates a user.
	Replace(ctx context.Context, userID string, userCtx flaggio.UserContext) error
	// Delete deletes a user.
	Delete(ctx context.Context, userID string) error
}
