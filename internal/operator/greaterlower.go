package operator

// Greater operator will check if the value from the user context is greater
// than any of the configured values on the flag.
func Greater(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := greater(v, usrValue, false)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// GreaterOrEqual operator will check if the value from the user context is greater
// or equal than any of the configured values on the flag.
func GreaterOrEqual(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := greater(v, usrValue, true)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// Lower operator will check if the value from the user context is lower
// than any of the configured values on the flag.
func Lower(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := lower(v, usrValue, false)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// LowerOrEqual operator will check if the value from the user context is lower
// or equal than any of the configured values on the flag.
func LowerOrEqual(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := lower(v, usrValue, true)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func greater(cnstrnValue, userValue interface{}, orEqual bool) (bool, error) {
	switch cv := cnstrnValue.(type) {
	case int, int32, int64:
		n1, _ := toInt64(cv)
		n2, err := toInt64(userValue)
		if err != nil {
			return false, err
		}
		if orEqual {
			return n2 >= n1, nil
		}
		return n2 > n1, nil
	case uint, uint32, uint64:
		n1, _ := toUInt64(cv)
		n2, err := toUInt64(userValue)
		if err != nil {
			return false, err
		}
		if orEqual {
			return n2 >= n1, nil
		}
		return n2 > n1, nil
	case float64:
		flt, ok := userValue.(float64)
		if !ok {
			return false, nil
		}
		if orEqual {
			return flt >= cv, nil
		}
		return flt > cv, nil
	default:
		return false, nil
	}
}

func lower(cnstrnValue, userValue interface{}, orEqual bool) (bool, error) {
	switch cv := cnstrnValue.(type) {
	case int, int32, int64:
		n1, _ := toInt64(cv)
		n2, err := toInt64(userValue)
		if err != nil {
			return false, err
		}
		if orEqual {
			return n2 <= n1, nil
		}
		return n2 < n1, nil
	case uint, uint32, uint64:
		n1, _ := toUInt64(cv)
		n2, err := toUInt64(userValue)
		if err != nil {
			return false, err
		}
		if orEqual {
			return n2 <= n1, nil
		}
		return n2 < n1, nil
	case float64:
		flt, ok := userValue.(float64)
		if !ok {
			return false, nil
		}
		if orEqual {
			return flt <= cv, nil
		}
		return flt < cv, nil
	default:
		return false, nil
	}
}
