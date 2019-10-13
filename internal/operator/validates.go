package operator

import (
	"github.com/sirupsen/logrus"
)

type Validator interface {
	Validate(usr map[string]interface{}) bool
}

type Validates struct{}

func (o Validates) Operate(usrValue interface{}, validValues []interface{}) bool {
	uv, ok := usrValue.(map[string]interface{})
	if !ok {
		logrus.WithField("value", uv).Error("expected user value to be a map[string]interface{}")
		return false
	}
	for _, v := range validValues {
		vldtr, ok := v.(Validator)
		if !ok || vldtr == nil {
			logrus.WithField("value", v).Error("expected value to be an Evaluator")
			return false
		}
		valid := vldtr.Validate(uv)
		if !valid {
			return false
		}
	}
	return true
}

type DoesntValidate struct{}

func (o DoesntValidate) Operate(usrValue interface{}, validValues []interface{}) bool {
	uv, ok := usrValue.(map[string]interface{})
	if !ok {
		logrus.WithField("value", uv).Error("expected user value to be a map[string]interface{}")
		return false
	}
	for _, v := range validValues {
		evltr, ok := v.(Validator)
		if !ok || evltr == nil {
			logrus.WithField("value", v).Error("expected value to be an Evaluator")
			return false
		}
		valid := evltr.Validate(uv)
		if valid {
			return false
		}
	}
	return true
}
