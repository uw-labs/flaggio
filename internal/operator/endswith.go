package operator

import (
	"strings"
)

// EndsWith operator will check if the value from the user context ends with
// any of the configured values on the flag.
func EndsWith(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		if endsWith(v, usrValue) {
			return true, nil
		}
	}
	return false, nil
}

// DoesntEndWith operator will check if the value from the user context doesn't end with
// any of the configured values on the flag.
func DoesntEndWith(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		if endsWith(v, usrValue) {
			return false, nil
		}
	}
	return true, nil
}

func endsWith(cnstrnValue, userValue interface{}) bool {
	str, err := toString(userValue)
	if err != nil {
		return false
	}
	switch v := cnstrnValue.(type) {
	case string:
		return strings.HasSuffix(str, v)
	case []byte:
		return strings.HasSuffix(str, string(v))
	default:
		return false
	}
}
