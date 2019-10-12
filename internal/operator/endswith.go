package operator

import (
	"strings"
)

type EndsWith struct{}

func (o EndsWith) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := endsWith(v, usrValue)
		if err != nil {
			return false
		}
		if ok {
			return true
		}
	}
	return false
}

type DoesntEndWith struct{}

func (o DoesntEndWith) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := endsWith(v, usrValue)
		if err != nil {
			return false
		}
		if ok {
			return false
		}
	}
	return true
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
