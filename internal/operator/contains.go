package operator

import (
	"errors"
	"strings"
)

// Contains operator will check if the value from the user context contains
// any of the configured values on the flag.
func Contains(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		if contains(v, usrValue) {
			return true, nil
		}
	}
	return false, nil
}

// DoesntContain operator will check if the value from the user context doesn't contain
// any of the configured values on the flag.
func DoesntContain(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		if contains(v, usrValue) {
			return false, nil
		}
	}
	return true, nil
}

func contains(cnstrnValue, userValue interface{}) bool {
	str, err := toString(userValue)
	if err != nil {
		return false
	}
	switch v := cnstrnValue.(type) {
	case string:
		return strings.Contains(str, v)
	case []byte:
		return strings.Contains(str, string(v))
	default:
		return false
	}
}

func toString(str interface{}) (string, error) {
	switch v := str.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", errors.New("invalid string type")
	}
}
