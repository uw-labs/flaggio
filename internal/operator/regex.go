package operator

import (
	"regexp"
)

// MatchesRegex operator will check if the value from the user context matches
// any regexes configured on the flag.
func MatchesRegex(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := matches(v, usrValue)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// DoesntMatchRegex operator will check if the value from the user context doesn't match
// any regexes configured on the flag.
func DoesntMatchRegex(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := matches(v, usrValue)
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
	}
	return true, nil
}

func matches(cnstrnValue, userValue interface{}) (bool, error) {
	str, err := toString(userValue)
	if err != nil {
		return false, nil
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
