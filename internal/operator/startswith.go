package operator

import (
	"strings"
)

// StartsWith operator will check if the value from the user context starts with
// any of the configured values on the flag.
func StartsWith(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		if startsWith(v, usrValue) {
			return true, nil
		}
	}
	return false, nil
}

// DoesntStartWith operator will check if the value from the user context doesn't start with
// any of the configured values on the flag.
func DoesntStartWith(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		if startsWith(v, usrValue) {
			return false, nil
		}
	}
	return true, nil
}

func startsWith(cnstrnValue, userValue interface{}) bool {
	str, err := toString(userValue)
	if err != nil {
		return false
	}
	switch v := cnstrnValue.(type) {
	case string:
		return strings.HasPrefix(str, v)
	case []byte:
		return strings.HasPrefix(str, string(v))
	default:
		return false
	}
}
