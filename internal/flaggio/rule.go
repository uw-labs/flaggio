package flaggio

type Rule struct {
	ID          string
	Constraints []*Constraint
}

func (r Rule) GetID() string {
	return r.ID
}

func (r Rule) Evaluate(usr map[string]interface{}) (EvalResult, error) {
	// op, ok := operator.FromProto(flaggio.Constraint_Operation(r.Operation))
	// if !ok {
	// 	// unknown operation, this is a configuration problem
	// 	logrus.WithField("operation", r.Operation).Error("unknown operation")
	// 	return evaluator.Result{}, errors.ErrUnknownOperator
	// }
	// if op.Operate(usr, r.Property, r.Values) {
	// 	return evaluator.Result{
	// 		Next: []evaluator.Evaluator{r.Distributions},
	// 	}, nil
	// }
	// TODO: fix
	return EvalResult{}, nil
}

type FlagRule struct {
	Rule
	Distributions []*Distribution
}

type SegmentRule struct {
	Rule
}
