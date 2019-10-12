package operator

import (
	"regexp"
)

type MatchesRegex struct{}

func (o MatchesRegex) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := matches(v, usrValue)
		if err != nil {
			return false
		}
		if ok {
			return true
		}
	}
	return false
}

type DoesntMatchRegex struct{}

func (o DoesntMatchRegex) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := matches(v, usrValue)
		if err != nil {
			return false
		}
		if ok {
			return false
		}
	}
	return true
}

func matches(cnstrnValue, userValue interface{}) (bool, error) {
	str, err := toString(userValue)
	if err != nil {
		return false, err
	}
	switch v := cnstrnValue.(type) {
	case string:
		return regexp.Match(v, []byte(str))
	case []byte:
		return regexp.Match(string(v), []byte(str))
	default:
		return false, nil
	}
}
