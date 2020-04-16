// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/victorkt/flaggio/internal/service (interfaces: Flag)

// Package service_mock is a generated GoMock package.
package service_mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	service "github.com/victorkt/flaggio/internal/service"
	reflect "reflect"
)

// MockFlag is a mock of Flag interface
type MockFlag struct {
	ctrl     *gomock.Controller
	recorder *MockFlagMockRecorder
}

// MockFlagMockRecorder is the mock recorder for MockFlag
type MockFlagMockRecorder struct {
	mock *MockFlag
}

// NewMockFlag creates a new mock instance
func NewMockFlag(ctrl *gomock.Controller) *MockFlag {
	mock := &MockFlag{ctrl: ctrl}
	mock.recorder = &MockFlagMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFlag) EXPECT() *MockFlagMockRecorder {
	return m.recorder
}

// Evaluate mocks base method
func (m *MockFlag) Evaluate(arg0 context.Context, arg1 string, arg2 *service.EvaluationRequest) (*service.EvaluationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Evaluate", arg0, arg1, arg2)
	ret0, _ := ret[0].(*service.EvaluationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Evaluate indicates an expected call of Evaluate
func (mr *MockFlagMockRecorder) Evaluate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Evaluate", reflect.TypeOf((*MockFlag)(nil).Evaluate), arg0, arg1, arg2)
}

// EvaluateAll mocks base method
func (m *MockFlag) EvaluateAll(arg0 context.Context, arg1 *service.EvaluationRequest) (*service.EvaluationsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EvaluateAll", arg0, arg1)
	ret0, _ := ret[0].(*service.EvaluationsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EvaluateAll indicates an expected call of EvaluateAll
func (mr *MockFlagMockRecorder) EvaluateAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EvaluateAll", reflect.TypeOf((*MockFlag)(nil).EvaluateAll), arg0, arg1)
}