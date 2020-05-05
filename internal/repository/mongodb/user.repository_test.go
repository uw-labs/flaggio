package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
)

func TestUserRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// drop database first
	if err := mongoDB.Drop(ctx); err != nil {
		t.Fatalf("failed drop database: %s", err)
	}

	// prepare users
	user1 := &flaggio.User{ID: "123", Context: flaggio.UserContext{"$userId": "123", "name": "john", "age": int32(33)}}
	user2 := &flaggio.User{ID: "456", Context: flaggio.UserContext{"$userId": "456", "name": "jane", "age": int32(46)}}

	// create new repo
	repo, err := mongo_repo.NewUserRepository(ctx, mongoDB)
	assert.NoError(t, err, "failed to create user repository")

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "create the first user",
			run: func(t *testing.T) {
				err = repo.Replace(ctx, user1.ID, user1.Context)
				assert.NoError(t, err, "failed to create first user")
			},
		},
		{
			name: "checks the user was created",
			run: func(t *testing.T) {
				usr, err := repo.FindByID(ctx, user1.ID)
				assert.NoError(t, err, "failed to find first user")
				user1.UpdatedAt = usr.UpdatedAt
				assert.Equal(t, user1, usr)
			},
		},
		{
			name: "create the second user",
			run: func(t *testing.T) {
				err = repo.Replace(ctx, user2.ID, user2.Context)
				assert.NoError(t, err, "failed to create second user")

			},
		},
		{
			name: "find the created user",
			run: func(t *testing.T) {
				usr, err := repo.FindByID(ctx, user2.ID)
				assert.NoError(t, err, "failed to find second user")
				user2.UpdatedAt = usr.UpdatedAt
				assert.Equal(t, user2, usr)
			},
		},
		{
			name: "find all users",
			run: func(t *testing.T) {
				usrs, err := repo.FindAll(ctx, nil, nil, nil)
				assert.NoError(t, err, "failed to find all users")
				assert.Equal(t, &flaggio.UserResults{Users: []*flaggio.User{user1, user2}, Total: 2}, usrs)
			},
		},
		{
			name: "search user",
			run: func(t *testing.T) {
				usrs, err := repo.FindAll(ctx, &user2.ID, nil, nil)
				assert.NoError(t, err, "failed to search user")
				assert.Equal(t, &flaggio.UserResults{Users: []*flaggio.User{user2}, Total: 1}, usrs)
			},
		},
		{
			name: "find all users with limit",
			run: func(t *testing.T) {
				usrs, err := repo.FindAll(ctx, nil, nil, int64Ptr(1))
				assert.NoError(t, err, "failed to all users with limit")
				assert.Equal(t, &flaggio.UserResults{Users: []*flaggio.User{user1}, Total: 2}, usrs)
			},
		},
		{
			name: "find all users with limit and offset",
			run: func(t *testing.T) {
				usrs, err := repo.FindAll(ctx, nil, int64Ptr(1), int64Ptr(1))
				assert.NoError(t, err, "failed to all users with limit and offset")
				assert.Equal(t, &flaggio.UserResults{Users: []*flaggio.User{user2}, Total: 2}, usrs)
			},
		},
		{
			name: "update the second user",
			run: func(t *testing.T) {
				user2.Context["validEmail"] = true
				err = repo.Replace(ctx, user2.ID, user2.Context)
				assert.NoError(t, err, "failed to update second user")
			},
		},
		{
			name: "find second user",
			run: func(t *testing.T) {
				usr, err := repo.FindByID(ctx, user2.ID)
				assert.NoError(t, err, "failed to find second user again")
				user2.UpdatedAt = usr.UpdatedAt
				assert.Equal(t, user2, usr)
			},
		},
		{
			name: "delete the first user",
			run: func(t *testing.T) {
				err = repo.Delete(ctx, user1.ID)
				assert.NoError(t, err, "failed to delete first user")

			},
		},
		{
			name: "find deleted user",
			run: func(t *testing.T) {
				usr, err := repo.FindByID(ctx, user1.ID)
				assert.EqualError(t, err, "user: not found")
				assert.Nil(t, usr)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, tt.run)
	}
}

func int64Ptr(n int64) *int64 {
	return &n
}
