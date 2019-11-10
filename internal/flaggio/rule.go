package flaggio

var _ Identifier = (*Rule)(nil)
var _ Evaluator = (*FlagRule)(nil)

// Rule has a list of constraints that all need to be satisfied so that
// it can pass.
type Rule struct {
	ID          string
	Constraints []*Constraint
}

// IsRuler is defined so that Rule can implement the Ruler interface.
func (r Rule) IsRuler() {}

// GetID returns the rule ID.
func (r Rule) GetID() string {
	return r.ID
}

// Populate will try to populate all references in the list of constraints.
func (r *Rule) Populate(identifiers []Identifier) {
	ConstraintList(r.Constraints).Populate(identifiers)
}

// FlagRule is a rule that also holds a list of distributions.
type FlagRule struct {
	Rule
	Distributions []*Distribution
}

// Evaluate will check that all constraints in this rule validates to true. If that
// is the case, it returns the list of distributions as next to be evaluated.
// If any of the constraints fail to pass, the rule returns an empty list of
// next evaluators. In any case, no answer is returned from the evaluation.
func (r FlagRule) Evaluate(usrContext map[string]interface{}) (EvalResult, error) {
	var next []Evaluator
	ok, err := ConstraintList(r.Constraints).Validate(usrContext)
	if ok {
		next = []Evaluator{DistributionList(r.Distributions)}
	}
	return EvalResult{
		Next: next,
	}, err
}

// SegmentRule is a rule to be used by segments.
type SegmentRule struct {
	Rule
}
