package operator

import (
	"strings"
)

type StartsWith struct{}

func (o StartsWith) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := startsWith(v, usrValue)
		if err != nil {
			return false
		}
		if ok {
			return true
		}
	}
	return false
}

type DoesntStartWith struct{}

func (o DoesntStartWith) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := startsWith(v, usrValue)
		if err != nil {
			return false
		}
		if ok {
			return false
		}
	}
	return true
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
