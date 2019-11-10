package operator

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/victorkt/flaggio/internal/errors"
)

// Validator validates the user context somehow and returns true when
// the validations are satisfied.
type Validator interface {
	Validate(usrContext map[string]interface{}) (bool, error)
}

// Validates operator will check if the value from the user context validates
// against the configured validator on the flag.
func Validates(usrValue interface{}, validValues []interface{}) (bool, error) {
	uv, ok := usrValue.(map[string]interface{})
	if !ok {
		return false, errors.InvalidFlag(
			fmt.Sprintf("expected user value to be a map[string]interface{}, got: %s", uv),
		)
	}
	for _, v := range validValues {
		vldtr, ok := v.(Validator)
		if !ok || vldtr == nil {
			// this happens when the reference is invalid on the flag
			// one possible scenario is the segment was deleted, leaving the
			// the flag constraint with an invalid reference
			logrus.WithField("value", v).Error("expected value to be an Evaluator")
			return false, nil
		}
		valid, err := vldtr.Validate(uv)
		if err != nil {
			return false, err
		}
		if !valid {
			return false, nil
		}
	}
	return true, nil
}

// DoesntValidate operator will check if the value from the user context doesn't validate
// against the configured validator on the flag.
func DoesntValidate(usrValue interface{}, validValues []interface{}) (bool, error) {
	uv, ok := usrValue.(map[string]interface{})
	if !ok {
		return false, errors.InvalidFlag(
			fmt.Sprintf("expected user value to be a map[string]interface{}, got: %s", uv),
		)
	}
	for _, v := range validValues {
		evltr, ok := v.(Validator)
		if !ok || evltr == nil {
			logrus.WithField("value", v).Error("expected value to be an Evaluator")
			return false, nil
		}
		valid, err := evltr.Validate(uv)
		if err != nil {
			return false, err
		}
		if valid {
			return false, nil
		}
	}
	return true, nil
}
