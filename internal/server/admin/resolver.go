package admin

import (
	"github.com/victorkt/flaggio/internal/repository"
)

type Resolver struct {
	flagRepo    repository.Flag
	variantRepo repository.Variant
	ruleRepo    repository.Rule
	segmentRepo repository.Segment
}

func NewResolver(
	flagRepo repository.Flag,
	variantRepo repository.Variant,
	ruleRepo repository.Rule,
	segmentRepo repository.Segment,
) *Resolver {
	return &Resolver{
		flagRepo:    flagRepo,
		variantRepo: variantRepo,
		ruleRepo:    ruleRepo,
		segmentRepo: segmentRepo,
	}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
