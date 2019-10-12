// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package flaggio

import (
	"fmt"
	"io"
	"strconv"
)

type Ruler interface {
	IsRuler()
}

type NewFlag struct {
	Key         string  `json:"key"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type NewVariant struct {
	Key            string      `json:"key"`
	Description    *string     `json:"description"`
	Value          interface{} `json:"value"`
	DefaultWhenOn  *bool       `json:"defaultWhenOn"`
	DefaultWhenOff *bool       `json:"defaultWhenOff"`
}

type Operation string

const (
	OperationOneOf            Operation = "ONE_OF"
	OperationNotOneOf         Operation = "NOT_ONE_OF"
	OperationGreater          Operation = "GREATER"
	OperationGreaterOrEqual   Operation = "GREATER_OR_EQUAL"
	OperationLower            Operation = "LOWER"
	OperationLowerOrEqual     Operation = "LOWER_OR_EQUAL"
	OperationExists           Operation = "EXISTS"
	OperationDoesntExist      Operation = "DOESNT_EXIST"
	OperationContains         Operation = "CONTAINS"
	OperationDoesntContain    Operation = "DOESNT_CONTAIN"
	OperationStartsWith       Operation = "STARTS_WITH"
	OperationDoesntStartWith  Operation = "DOESNT_START_WITH"
	OperationEndsWith         Operation = "ENDS_WITH"
	OperationDoesntEndWith    Operation = "DOESNT_END_WITH"
	OperationMatchesRegex     Operation = "MATCHES_REGEX"
	OperationDoesntMatchRegex Operation = "DOESNT_MATCH_REGEX"
	OperationBeforeDate       Operation = "BEFORE_DATE"
	OperationBeforeOrSameDate Operation = "BEFORE_OR_SAME_DATE"
	OperationAfterDate        Operation = "AFTER_DATE"
	OperationAfterOrSameDate  Operation = "AFTER_OR_SAME_DATE"
	OperationIsInSegment      Operation = "IS_IN_SEGMENT"
	OperationIsntInSegment    Operation = "ISNT_IN_SEGMENT"
)

var AllOperation = []Operation{
	OperationOneOf,
	OperationNotOneOf,
	OperationGreater,
	OperationGreaterOrEqual,
	OperationLower,
	OperationLowerOrEqual,
	OperationExists,
	OperationDoesntExist,
	OperationContains,
	OperationDoesntContain,
	OperationStartsWith,
	OperationDoesntStartWith,
	OperationEndsWith,
	OperationDoesntEndWith,
	OperationMatchesRegex,
	OperationDoesntMatchRegex,
	OperationBeforeDate,
	OperationBeforeOrSameDate,
	OperationAfterDate,
	OperationAfterOrSameDate,
	OperationIsInSegment,
	OperationIsntInSegment,
}

func (e Operation) IsValid() bool {
	switch e {
	case OperationOneOf, OperationNotOneOf, OperationGreater, OperationGreaterOrEqual, OperationLower, OperationLowerOrEqual, OperationExists, OperationDoesntExist, OperationContains, OperationDoesntContain, OperationStartsWith, OperationDoesntStartWith, OperationEndsWith, OperationDoesntEndWith, OperationMatchesRegex, OperationDoesntMatchRegex, OperationBeforeDate, OperationBeforeOrSameDate, OperationAfterDate, OperationAfterOrSameDate, OperationIsInSegment, OperationIsntInSegment:
		return true
	}
	return false
}

func (e Operation) String() string {
	return string(e)
}

func (e *Operation) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Operation(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Operation", str)
	}
	return nil
}

func (e Operation) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
