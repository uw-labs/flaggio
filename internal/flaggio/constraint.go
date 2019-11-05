package flaggio

//go:generate mockgen -destination=./mocks/operator_mock.go -package=flaggio_mock github.com/victorkt/flaggio/internal/flaggio Operator

import (
	"github.com/sirupsen/logrus"
	"github.com/victorkt/flaggio/internal/operator"
)

var _ operator.Validator = (*Constraint)(nil)

type Operator interface {
	Operate(usrValue interface{}, validValues []interface{}) bool
}

type Constraint struct {
	ID        string
	Property  string
	Operation Operation
	Values    []interface{}
}

func (c Constraint) Validate(usrContext map[string]interface{}) bool {
	op, ok := operatorMap[c.Operation]
	if !ok {
		// unknown operation, this is a configuration problem
		logrus.WithField("operation", c.Operation).Error("unknown operation")
		return false
	}
	// TODO: check if worth to return errors
	switch c.Operation {
	case OperationIsInSegment, OperationIsntInSegment:
		return op.Operate(usrContext, c.Values)
	default:
		return op.Operate(usrContext[c.Property], c.Values)
	}
}

func (c *Constraint) Populate(identifiers []Identifier) {
	switch c.Operation {
	case OperationIsInSegment, OperationIsntInSegment:
		c.populateSegments(identifiers)
	}
}

func (c *Constraint) populateSegments(identifiers []Identifier) {
	for idx := 0; idx < len(c.Values); idx++ {
		id := c.Values[idx]
		for _, ider := range identifiers {
			if _, ok := ider.(*Segment); ok && ider.GetID() == id {
				c.Values[idx] = ider
			}
		}
	}
}

type ConstraintList []*Constraint

func (l ConstraintList) Validate(usrContext map[string]interface{}) bool {
	for _, c := range l {
		if !c.Validate(usrContext) {
			return false
		}
	}
	return true
}

func (l ConstraintList) Populate(identifiers []Identifier) {
	for _, c := range l {
		c.Populate(identifiers)
	}
}

var operatorMap = map[Operation]Operator{
	OperationOneOf:            operator.OneOf{},
	OperationNotOneOf:         operator.NotOneOf{},
	OperationGreater:          operator.Greater{},
	OperationGreaterOrEqual:   operator.GreaterOrEqual{},
	OperationLower:            operator.Lower{},
	OperationLowerOrEqual:     operator.LowerOrEqual{},
	OperationExists:           operator.Exists{},
	OperationDoesntExist:      operator.DoesntExist{},
	OperationContains:         operator.Contains{},
	OperationDoesntContain:    operator.DoesntContain{},
	OperationStartsWith:       operator.StartsWith{},
	OperationDoesntStartWith:  operator.DoesntStartWith{},
	OperationEndsWith:         operator.EndsWith{},
	OperationDoesntEndWith:    operator.DoesntEndWith{},
	OperationMatchesRegex:     operator.MatchesRegex{},
	OperationDoesntMatchRegex: operator.DoesntMatchRegex{},
	// TODO: date operations
	// TODO: ip address operator
	OperationIsInSegment:   operator.Validates{},
	OperationIsntInSegment: operator.DoesntValidate{},
}
