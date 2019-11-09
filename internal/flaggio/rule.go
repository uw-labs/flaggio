package flaggio

var _ Identifier = (*Rule)(nil)
var _ Evaluator = (*FlagRule)(nil)

type Rule struct {
	ID          string
	Constraints []*Constraint
}

func (r Rule) IsRuler() {}

func (r Rule) GetID() string {
	return r.ID
}

func (r *Rule) Populate(identifiers []Identifier) {
	ConstraintList(r.Constraints).Populate(identifiers)
}

type FlagRule struct {
	Rule
	Distributions []*Distribution
}

func (r FlagRule) Evaluate(usrContext map[string]interface{}) (EvalResult, error) {
	var next []Evaluator
	if ConstraintList(r.Constraints).Validate(usrContext) {
		next = []Evaluator{DistributionList(r.Distributions)}
	}
	return EvalResult{
		Next: next,
	}, nil
}

type SegmentRule struct {
	Rule
}
