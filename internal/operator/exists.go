package operator

// Exists operator will check if the property from the user context exists (is not null)
func Exists(usrValue interface{}, _ []interface{}) (bool, error) {
	return usrValue != nil, nil
}

// DoesntExist operator will check if the property from the user context doesn't exist (is null)
func DoesntExist(usrValue interface{}, _ []interface{}) (bool, error) {
	return usrValue == nil, nil
}
