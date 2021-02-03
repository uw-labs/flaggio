// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/uw-labs/flaggio/internal/repository (interfaces: Evaluation)

// Package repository_mock is a generated GoMock package.
package repository_mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	flaggio "github.com/uw-labs/flaggio/internal/flaggio"
	reflect "reflect"
)

// MockEvaluation is a mock of Evaluation interface
type MockEvaluation struct {
	ctrl     *gomock.Controller
	recorder *MockEvaluationMockRecorder
}

// MockEvaluationMockRecorder is the mock recorder for MockEvaluation
type MockEvaluationMockRecorder struct {
	mock *MockEvaluation
}

// NewMockEvaluation creates a new mock instance
func NewMockEvaluation(ctrl *gomock.Controller) *MockEvaluation {
	mock := &MockEvaluation{ctrl: ctrl}
	mock.recorder = &MockEvaluationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEvaluation) EXPECT() *MockEvaluationMockRecorder {
	return m.recorder
}

// DeleteAllByUserID mocks base method
func (m *MockEvaluation) DeleteAllByUserID(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllByUserID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllByUserID indicates an expected call of DeleteAllByUserID
func (mr *MockEvaluationMockRecorder) DeleteAllByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllByUserID", reflect.TypeOf((*MockEvaluation)(nil).DeleteAllByUserID), arg0, arg1)
}

// DeleteByID mocks base method
func (m *MockEvaluation) DeleteByID(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID
func (mr *MockEvaluationMockRecorder) DeleteByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockEvaluation)(nil).DeleteByID), arg0, arg1)
}

// FindAllByReqHash mocks base method
func (m *MockEvaluation) FindAllByReqHash(arg0 context.Context, arg1 string) (flaggio.EvaluationList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByReqHash", arg0, arg1)
	ret0, _ := ret[0].(flaggio.EvaluationList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByReqHash indicates an expected call of FindAllByReqHash
func (mr *MockEvaluationMockRecorder) FindAllByReqHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByReqHash", reflect.TypeOf((*MockEvaluation)(nil).FindAllByReqHash), arg0, arg1)
}

// FindAllByUserID mocks base method
func (m *MockEvaluation) FindAllByUserID(arg0 context.Context, arg1 string, arg2 *string, arg3, arg4 *int64) (*flaggio.EvaluationResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByUserID", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*flaggio.EvaluationResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByUserID indicates an expected call of FindAllByUserID
func (mr *MockEvaluationMockRecorder) FindAllByUserID(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByUserID", reflect.TypeOf((*MockEvaluation)(nil).FindAllByUserID), arg0, arg1, arg2, arg3, arg4)
}

// FindByID mocks base method
func (m *MockEvaluation) FindByID(arg0 context.Context, arg1 string) (*flaggio.Evaluation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*flaggio.Evaluation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockEvaluationMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockEvaluation)(nil).FindByID), arg0, arg1)
}

// FindByReqHashAndFlagKey mocks base method
func (m *MockEvaluation) FindByReqHashAndFlagKey(arg0 context.Context, arg1, arg2 string) (*flaggio.Evaluation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByReqHashAndFlagKey", arg0, arg1, arg2)
	ret0, _ := ret[0].(*flaggio.Evaluation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByReqHashAndFlagKey indicates an expected call of FindByReqHashAndFlagKey
func (mr *MockEvaluationMockRecorder) FindByReqHashAndFlagKey(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByReqHashAndFlagKey", reflect.TypeOf((*MockEvaluation)(nil).FindByReqHashAndFlagKey), arg0, arg1, arg2)
}

// ReplaceAll mocks base method
func (m *MockEvaluation) ReplaceAll(arg0 context.Context, arg1, arg2 string, arg3 flaggio.EvaluationList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceAll", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceAll indicates an expected call of ReplaceAll
func (mr *MockEvaluationMockRecorder) ReplaceAll(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceAll", reflect.TypeOf((*MockEvaluation)(nil).ReplaceAll), arg0, arg1, arg2, arg3)
}

// ReplaceOne mocks base method
func (m *MockEvaluation) ReplaceOne(arg0 context.Context, arg1 string, arg2 *flaggio.Evaluation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceOne", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceOne indicates an expected call of ReplaceOne
func (mr *MockEvaluationMockRecorder) ReplaceOne(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceOne", reflect.TypeOf((*MockEvaluation)(nil).ReplaceOne), arg0, arg1, arg2)
}
