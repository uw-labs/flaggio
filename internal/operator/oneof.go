package operator

import (
	"encoding/json"
	"errors"
)

type OneOf struct{}

func (o OneOf) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := equals(v, usrValue)
		if err != nil {
			return false
		}
		if ok {
			return true
		}
	}
	return false
}

type NotOneOf struct{}

func (o NotOneOf) Operate(usrValue interface{}, validValues []interface{}) bool {
	for _, v := range validValues {
		ok, err := equals(v, usrValue)
		if err != nil || ok {
			return false
		}
	}
	return true
}

func equals(cnstrnValue, userValue interface{}) (bool, error) {
	switch v := cnstrnValue.(type) {
	// float32 <-> float64 conversions return different number precisions
	case string, bool, float64:
		return v == userValue, nil
	case int, int32, int64:
		v1, _ := toInt64(v)
		v2, err := toInt64(userValue)
		if err != nil {
			return false, err
		}
		return v1 == v2, nil
	case uint, uint32, uint64:
		v1, _ := toUInt64(v)
		v2, err := toUInt64(userValue)
		if err != nil {
			return false, err
		}
		return v1 == v2, nil
	case []byte:
		uv, ok := userValue.([]byte)
		if !ok {
			return false, nil
		}
		return string(v) == string(uv), nil
	default:
		return false, nil
	}
}

func toInt64(n interface{}) (int64, error) {
	switch v := n.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case json.Number:
		return v.Int64()
	default:
		return 0, errors.New("not an integer")
	}
}

func toUInt64(n interface{}) (uint64, error) {
	switch v := n.(type) {
	case uint:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return v, nil
	default:
		return 0, errors.New("not an unsigned integer")
	}
}
