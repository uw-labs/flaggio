package operator

// Exists operator will check if the property from the user context exists (is not null)
type Exists struct{}

// Operate will check the value from the user context against the configured validators
func (o Exists) Operate(usrValue interface{}, _ []interface{}) (bool, error) {
	return usrValue != nil, nil
}

// DoesntExist operator will check if the property from the user context doesn't exist (is null)
type DoesntExist struct{}

// Operate will check the value from the user context against the configured validators
func (o DoesntExist) Operate(usrValue interface{}, _ []interface{}) (bool, error) {
	return usrValue == nil, nil
}
