package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrNotImplemented        = errors.New("not implemented")
	ErrNoDefaultVariant      = errors.New("no default variant defined for flag")
	ErrNoVariantToDistribute = errors.New("no variants to distribute, please check the rule configuration")
	ErrUnknownOperator       = errors.New("unknown operator, please check the rule configuration")
	ErrNothingToUpdate       = errors.New("nothing to update, please specify fields to update")
)

func NotFound(entity string) error {
	return errors.Errorf("%s not found", entity)
}
