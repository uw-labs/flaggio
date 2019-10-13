package api

import (
	"context"
)

var _ MutationResolver = &mutationResolver{}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Ping(ctx context.Context) (*bool, error) {
	pong := true
	return &pong, nil
}
