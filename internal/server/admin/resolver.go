package admin

import (
	"github.com/victorkt/flaggio/internal/repository"
)

// Resolver is the root resolver for the GraphQL server.
type Resolver struct {
	FlagRepo    repository.Flag
	VariantRepo repository.Variant
	RuleRepo    repository.Rule
	SegmentRepo repository.Segment
}

// Mutation returns the mutation resolver.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns the query resolver.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
