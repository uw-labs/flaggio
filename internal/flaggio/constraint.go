package flaggio

type Constraint struct {
	ID        string        `json:"id"`
	Property  string        `json:"property"`
	Operation Operation     `json:"operation"`
	Values    []interface{} `json:"values"`
}

func (c Constraint) Evaluate(usr map[string]interface{}) (bool, error) {
	// op, ok := operator.FromProto(flaggio.Constraint_Operation(c.Operation))
	// if !ok {
	// 	// unknown operation, this is a configuration problem
	// 	logrus.WithField("operation", c.Operation).Error("unknown operation")
	// 	return evaluator.Result{}, errors.ErrUnknownOperator
	// }
	// if op.Operate(usr, c.Property, c.Values) {
	// 	// TODO: fix
	// 	// return evaluator.Result{
	// 	// 	Next: []evaluator.Evaluator{c.Distributions},
	// 	// }, nil
	// }
	// return evaluator.Result{}, nil
	return true, nil
}

type ConstraintList []*Constraint

func (l ConstraintList) Evaluate(usr map[string]interface{}) (bool, error) {
	for _, c := range l {
		pass, err := c.Evaluate(usr)
		if err != nil {
			return false, err
		}
		if pass {
			return true, nil
		}
	}
	return false, nil
}
