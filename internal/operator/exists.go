package operator

type Exists struct{}

func (o Exists) Operate(usrValue interface{}, _ []interface{}) bool {
	return usrValue != nil
}

type DoesntExist struct{}

func (o DoesntExist) Operate(usrValue interface{}, _ []interface{}) bool {
	return usrValue == nil
}
