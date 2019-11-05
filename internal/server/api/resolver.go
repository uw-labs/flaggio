package api

import (
	"github.com/victorkt/flaggio/internal/repository"
)

type Resolver struct {
	flagRepo    repository.Flag
	segmentRepo repository.Segment
}

func NewResolver(
	flagRepo repository.Flag,
	segmentRepo repository.Segment,
) *Resolver {
	return &Resolver{
		flagRepo:    flagRepo,
		segmentRepo: segmentRepo,
	}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
