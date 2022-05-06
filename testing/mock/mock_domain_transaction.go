// Code generated by MockGen. DO NOT EDIT.
// Source: transaction.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/goropikari/simpledbgo/domain"
)

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// BlockLength mocks base method.
func (m *MockTransaction) BlockLength(arg0 domain.FileName) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockLength", arg0)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockLength indicates an expected call of BlockLength.
func (mr *MockTransactionMockRecorder) BlockLength(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockLength", reflect.TypeOf((*MockTransaction)(nil).BlockLength), arg0)
}

// BlockSize mocks base method.
func (m *MockTransaction) BlockSize() domain.BlockSize {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockSize")
	ret0, _ := ret[0].(domain.BlockSize)
	return ret0
}

// BlockSize indicates an expected call of BlockSize.
func (mr *MockTransactionMockRecorder) BlockSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockSize", reflect.TypeOf((*MockTransaction)(nil).BlockSize))
}

// Commit mocks base method.
func (m *MockTransaction) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTransactionMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTransaction)(nil).Commit))
}

// ExtendFile mocks base method.
func (m *MockTransaction) ExtendFile(arg0 domain.FileName) (domain.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExtendFile", arg0)
	ret0, _ := ret[0].(domain.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExtendFile indicates an expected call of ExtendFile.
func (mr *MockTransactionMockRecorder) ExtendFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtendFile", reflect.TypeOf((*MockTransaction)(nil).ExtendFile), arg0)
}

// GetInt32 mocks base method.
func (m *MockTransaction) GetInt32(arg0 domain.Block, arg1 int64) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt32", arg0, arg1)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInt32 indicates an expected call of GetInt32.
func (mr *MockTransactionMockRecorder) GetInt32(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt32", reflect.TypeOf((*MockTransaction)(nil).GetInt32), arg0, arg1)
}

// GetString mocks base method.
func (m *MockTransaction) GetString(arg0 domain.Block, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetString indicates an expected call of GetString.
func (mr *MockTransactionMockRecorder) GetString(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockTransaction)(nil).GetString), arg0, arg1)
}

// Pin mocks base method.
func (m *MockTransaction) Pin(arg0 domain.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pin", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Pin indicates an expected call of Pin.
func (mr *MockTransactionMockRecorder) Pin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pin", reflect.TypeOf((*MockTransaction)(nil).Pin), arg0)
}

// SetInt32 mocks base method.
func (m *MockTransaction) SetInt32(arg0 domain.Block, arg1 int64, arg2 int32, arg3 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetInt32", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInt32 indicates an expected call of SetInt32.
func (mr *MockTransactionMockRecorder) SetInt32(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInt32", reflect.TypeOf((*MockTransaction)(nil).SetInt32), arg0, arg1, arg2, arg3)
}

// SetString mocks base method.
func (m *MockTransaction) SetString(arg0 domain.Block, arg1 int64, arg2 string, arg3 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetString", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetString indicates an expected call of SetString.
func (mr *MockTransactionMockRecorder) SetString(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetString", reflect.TypeOf((*MockTransaction)(nil).SetString), arg0, arg1, arg2, arg3)
}

// Unpin mocks base method.
func (m *MockTransaction) Unpin(arg0 domain.Block) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unpin", arg0)
}

// Unpin indicates an expected call of Unpin.
func (mr *MockTransactionMockRecorder) Unpin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unpin", reflect.TypeOf((*MockTransaction)(nil).Unpin), arg0)
}
