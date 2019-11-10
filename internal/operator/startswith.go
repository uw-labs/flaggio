package operator

import (
	"strings"
)

// StartsWith operator will check if the value from the user context starts with
// any of the configured values on the flag.
func StartsWith(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := startsWith(v, usrValue)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// DoesntStartWith operator will check if the value from the user context doesn't start with
// any of the configured values on the flag.
func DoesntStartWith(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := startsWith(v, usrValue)
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
	}
	return true, nil
}

func startsWith(cnstrnValue, userValue interface{}) (bool, error) {
	str, err := toString(userValue)
	if err != nil {
		return false, err
	}
	switch v := cnstrnValue.(type) {
	case string:
		return strings.HasPrefix(str, v), nil
	case []byte:
		return strings.HasPrefix(str, string(v)), nil
	default:
		return false, nil
	}
}
