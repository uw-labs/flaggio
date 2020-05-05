package admin

import (
	"github.com/victorkt/flaggio/internal/repository"
)

var _ ResolverRoot = (*Resolver)(nil)

// Resolver is the root resolver for the GraphQL server.
type Resolver struct {
	FlagRepo       repository.Flag
	VariantRepo    repository.Variant
	RuleRepo       repository.Rule
	SegmentRepo    repository.Segment
	UserRepo       repository.User
	EvaluationRepo repository.Evaluation
}

// Mutation returns the mutation resolver.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns the query resolver.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// User returns the user resolver.
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
