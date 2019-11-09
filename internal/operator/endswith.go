package operator

import (
	"strings"
)

// EndsWith operator will check if the value from the user context ends with
// any of the configured values on the flag.
type EndsWith struct{}

// Operate will check the value from the user context against the configured validators
func (o EndsWith) Operate(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := endsWith(v, usrValue)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// DoesntEndWith operator will check if the value from the user context doesn't end with
// any of the configured values on the flag.
type DoesntEndWith struct{}

// Operate will check the value from the user context against the configured validators
func (o DoesntEndWith) Operate(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := endsWith(v, usrValue)
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
	}
	return true, nil
}

func endsWith(cnstrnValue, userValue interface{}) (bool, error) {
	str, err := toString(userValue)
	if err != nil {
		return false, err
	}
	switch v := cnstrnValue.(type) {
	case string:
		return strings.HasSuffix(str, v), nil
	case []byte:
		return strings.HasSuffix(str, string(v)), nil
	default:
		return false, nil
	}
}
