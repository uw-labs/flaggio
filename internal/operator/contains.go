package operator

import (
	"errors"
	"strings"
)

// Contains operator will check if the value from the user context contains
// any of the configured values on the flag.
func Contains(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := contains(v, usrValue)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// DoesntContain operator will check if the value from the user context doesn't contain
// any of the configured values on the flag.
func DoesntContain(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := contains(v, usrValue)
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
	}
	return true, nil
}

func contains(cnstrnValue, userValue interface{}) (bool, error) {
	str, err := toString(userValue)
	if err != nil {
		return false, err
	}
	switch v := cnstrnValue.(type) {
	case string:
		return strings.Contains(str, v), nil
	case []byte:
		return strings.Contains(str, string(v)), nil
	default:
		return false, nil
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
