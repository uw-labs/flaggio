package operator

import (
	"errors"
	"strings"
)

type Contains struct{}

func (o Contains) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := contains(v, usrValue)
		if err != nil {
			return false
		}
		if ok {
			return true
		}
	}
	return false
}

type DoesntContain struct{}

func (o DoesntContain) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := contains(v, usrValue)
		if err != nil {
			return false
		}
		if ok {
			return false
		}
	}
	return true
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
