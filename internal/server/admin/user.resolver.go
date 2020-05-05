package admin

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

var _ UserResolver = &userResolver{}

type userResolver struct{ *Resolver }

// Evaluations returns the list of evaluations for a given user.
func (r *userResolver) Evaluations(ctx context.Context, usr *flaggio.User, search *string, offset, limit *int) (*flaggio.EvaluationResults, error) {
	var ofst, lmt *int64
	if offset != nil {
		v := int64(*offset)
		ofst = &v
	}
	if limit != nil {
		v := int64(*limit)
		lmt = &v
	}
	return r.EvaluationRepo.FindAllByUserID(ctx, usr.ID, search, ofst, lmt)
}
