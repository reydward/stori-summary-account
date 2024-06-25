// Code generated by MockGen. DO NOT EDIT.
// Source: summary/internal/repository (interfaces: SummaryRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	model "summary/internal/model"
)

// MockSummaryRepository is a mock of SummaryRepository interface.
type MockSummaryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSummaryRepositoryMockRecorder
}

// MockSummaryRepositoryMockRecorder is the mock recorder for MockSummaryRepository.
type MockSummaryRepositoryMockRecorder struct {
	mock *MockSummaryRepository
}

// NewMockSummaryRepository creates a new mock instance.
func NewMockSummaryRepository(ctrl *gomock.Controller) *MockSummaryRepository {
	mock := &MockSummaryRepository{ctrl: ctrl}
	mock.recorder = &MockSummaryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSummaryRepository) EXPECT() *MockSummaryRepositoryMockRecorder {
	return m.recorder
}

// GetAccountAverageCredit mocks base method.
func (m *MockSummaryRepository) GetAccountAverageCredit(arg0 int) (float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountAverageCredit", arg0)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountAverageCredit indicates an expected call of GetAccountAverageCredit.
func (mr *MockSummaryRepositoryMockRecorder) GetAccountAverageCredit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountAverageCredit", reflect.TypeOf((*MockSummaryRepository)(nil).GetAccountAverageCredit), arg0)
}

// GetAccountAverageDebit mocks base method.
func (m *MockSummaryRepository) GetAccountAverageDebit(arg0 int) (float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountAverageDebit", arg0)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountAverageDebit indicates an expected call of GetAccountAverageDebit.
func (mr *MockSummaryRepositoryMockRecorder) GetAccountAverageDebit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountAverageDebit", reflect.TypeOf((*MockSummaryRepository)(nil).GetAccountAverageDebit), arg0)
}

// GetAccountInfo mocks base method.
func (m *MockSummaryRepository) GetAccountInfo(arg0 int) (*model.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountInfo", arg0)
	ret0, _ := ret[0].(*model.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountInfo indicates an expected call of GetAccountInfo.
func (mr *MockSummaryRepositoryMockRecorder) GetAccountInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountInfo", reflect.TypeOf((*MockSummaryRepository)(nil).GetAccountInfo), arg0)
}

// GetAccountNumberOfTransactions mocks base method.
func (m *MockSummaryRepository) GetAccountNumberOfTransactions(arg0 int) ([]*model.NumberOfTransactions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountNumberOfTransactions", arg0)
	ret0, _ := ret[0].([]*model.NumberOfTransactions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountNumberOfTransactions indicates an expected call of GetAccountNumberOfTransactions.
func (mr *MockSummaryRepositoryMockRecorder) GetAccountNumberOfTransactions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountNumberOfTransactions", reflect.TypeOf((*MockSummaryRepository)(nil).GetAccountNumberOfTransactions), arg0)
}

// GetAccountTotalBalance mocks base method.
func (m *MockSummaryRepository) GetAccountTotalBalance(arg0 int) (float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountTotalBalance", arg0)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountTotalBalance indicates an expected call of GetAccountTotalBalance.
func (mr *MockSummaryRepositoryMockRecorder) GetAccountTotalBalance(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountTotalBalance", reflect.TypeOf((*MockSummaryRepository)(nil).GetAccountTotalBalance), arg0)
}

// GetUser mocks base method.
func (m *MockSummaryRepository) GetUser(arg0 int) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockSummaryRepositoryMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockSummaryRepository)(nil).GetUser), arg0)
}
