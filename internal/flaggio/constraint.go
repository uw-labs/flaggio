package flaggio

import (
	"fmt"

	"github.com/victorkt/flaggio/internal/errors"
	"github.com/victorkt/flaggio/internal/operator"
)

var _ operator.Validator = (*Constraint)(nil)

// Operator runs some logic using the value from the user context and the
// value from the flag configuration to determine if the operation is
// satisfied or not.
type Operator func(usrValue interface{}, validValues []interface{}) (bool, error)

// Constraint holds all the information needed for an operation to be executed.
type Constraint struct {
	ID        string
	Property  string
	Operation Operation
	Values    []interface{}
}

// Validate will check if a property in the user context passes some operation based on
// some configured valid values. Some operations don't need the property to be defined,
// while some others don't required any valid values.
func (c Constraint) Validate(usrContext map[string]interface{}) (bool, error) {
	operate, ok := operatorMap[c.Operation]
	if !ok {
		// unknown operation, this is a configuration problem
		return false, errors.InvalidFlag(fmt.Sprintf("unknown operation: %s", c.Operation))
	}
	switch c.Operation {
	case OperationIsInSegment, OperationIsntInSegment:
		return operate(usrContext, c.Values)
	default:
		return operate(usrContext[c.Property], c.Values)
	}
}

// Populate will try to populate all references on this constraint
// depending on the operation type. For OperationIsInSegment and
// OperationIsntInSegment, the flag will have a reference to a segment.
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

// ConstraintList is a slice of *Constraint
type ConstraintList []*Constraint

// Validate will check if all the constraints on this list passes their validation.
// If any of the constraints don't validate to true, the result is false.
func (l ConstraintList) Validate(usrContext map[string]interface{}) (bool, error) {
	for _, c := range l {
		ok, err := c.Validate(usrContext)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// Populate will try to populate all references on this constraint list
func (l ConstraintList) Populate(identifiers []Identifier) {
	for _, c := range l {
		c.Populate(identifiers)
	}
}

// Maps the GraphQL enum to the operator func
var operatorMap = map[Operation]Operator{
	OperationOneOf:            operator.OneOf,
	OperationNotOneOf:         operator.NotOneOf,
	OperationGreater:          operator.Greater,
	OperationGreaterOrEqual:   operator.GreaterOrEqual,
	OperationLower:            operator.Lower,
	OperationLowerOrEqual:     operator.LowerOrEqual,
	OperationExists:           operator.Exists,
	OperationDoesntExist:      operator.DoesntExist,
	OperationContains:         operator.Contains,
	OperationDoesntContain:    operator.DoesntContain,
	OperationStartsWith:       operator.StartsWith,
	OperationDoesntStartWith:  operator.DoesntStartWith,
	OperationEndsWith:         operator.EndsWith,
	OperationDoesntEndWith:    operator.DoesntEndWith,
	OperationMatchesRegex:     operator.MatchesRegex,
	OperationDoesntMatchRegex: operator.DoesntMatchRegex,
	OperationIsInSegment:      operator.Validates,
	OperationIsntInSegment:    operator.DoesntValidate,
}
