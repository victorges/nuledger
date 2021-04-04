// Code generated by MockGen. DO NOT EDIT.
// Source: ./ledger.go

// Package mocks is a generated GoMock package.
package mocks

import (
	model "nuledger/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLedger is a mock of Ledger interface.
type MockLedger struct {
	ctrl     *gomock.Controller
	recorder *MockLedgerMockRecorder
}

// MockLedgerMockRecorder is the mock recorder for MockLedger.
type MockLedgerMockRecorder struct {
	mock *MockLedger
}

// NewMockLedger creates a new mock instance.
func NewMockLedger(ctrl *gomock.Controller) *MockLedger {
	mock := &MockLedger{ctrl: ctrl}
	mock.recorder = &MockLedgerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLedger) EXPECT() *MockLedgerMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockLedger) CreateAccount(account model.Account) (*model.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", account)
	ret0, _ := ret[0].(*model.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockLedgerMockRecorder) CreateAccount(account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockLedger)(nil).CreateAccount), account)
}

// PerformTransaction mocks base method.
func (m *MockLedger) PerformTransaction(transaction model.Transaction) (*model.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PerformTransaction", transaction)
	ret0, _ := ret[0].(*model.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PerformTransaction indicates an expected call of PerformTransaction.
func (mr *MockLedgerMockRecorder) PerformTransaction(transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PerformTransaction", reflect.TypeOf((*MockLedger)(nil).PerformTransaction), transaction)
}