// Code generated by MockGen. DO NOT EDIT.
// Source: ./authorizer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	rule "nuledger/authorizer/rule"
	model "nuledger/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthorizer is a mock of Authorizer interface.
type MockAuthorizer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizerMockRecorder
}

// MockAuthorizerMockRecorder is the mock recorder for MockAuthorizer.
type MockAuthorizerMockRecorder struct {
	mock *MockAuthorizer
}

// NewMockAuthorizer creates a new mock instance.
func NewMockAuthorizer(ctrl *gomock.Controller) *MockAuthorizer {
	mock := &MockAuthorizer{ctrl: ctrl}
	mock.recorder = &MockAuthorizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizer) EXPECT() *MockAuthorizerMockRecorder {
	return m.recorder
}

// Authorize mocks base method.
func (m *MockAuthorizer) Authorize(account model.Account, transaction model.Transaction) (rule.CommitFunc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", account, transaction)
	ret0, _ := ret[0].(rule.CommitFunc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorize indicates an expected call of Authorize.
func (mr *MockAuthorizerMockRecorder) Authorize(account, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockAuthorizer)(nil).Authorize), account, transaction)
}