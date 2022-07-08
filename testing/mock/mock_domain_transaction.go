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
func (m *MockTransaction) GetInt32(blk domain.Block, offset int64) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt32", blk, offset)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInt32 indicates an expected call of GetInt32.
func (mr *MockTransactionMockRecorder) GetInt32(blk, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt32", reflect.TypeOf((*MockTransaction)(nil).GetInt32), blk, offset)
}

// GetString mocks base method.
func (m *MockTransaction) GetString(blk domain.Block, offset int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", blk, offset)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetString indicates an expected call of GetString.
func (mr *MockTransactionMockRecorder) GetString(blk, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockTransaction)(nil).GetString), blk, offset)
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

// Recover mocks base method.
func (m *MockTransaction) Recover() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recover")
	ret0, _ := ret[0].(error)
	return ret0
}

// Recover indicates an expected call of Recover.
func (mr *MockTransactionMockRecorder) Recover() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recover", reflect.TypeOf((*MockTransaction)(nil).Recover))
}

// Rollback mocks base method.
func (m *MockTransaction) Rollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTransactionMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTransaction)(nil).Rollback))
}

// SetInt32 mocks base method.
func (m *MockTransaction) SetInt32(blk domain.Block, offset int64, val int32, writeLog bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetInt32", blk, offset, val, writeLog)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInt32 indicates an expected call of SetInt32.
func (mr *MockTransactionMockRecorder) SetInt32(blk, offset, val, writeLog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInt32", reflect.TypeOf((*MockTransaction)(nil).SetInt32), blk, offset, val, writeLog)
}

// SetString mocks base method.
func (m *MockTransaction) SetString(blk domain.Block, offset int64, val string, writeLog bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetString", blk, offset, val, writeLog)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetString indicates an expected call of SetString.
func (mr *MockTransactionMockRecorder) SetString(blk, offset, val, writeLog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetString", reflect.TypeOf((*MockTransaction)(nil).SetString), blk, offset, val, writeLog)
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

// MockTxNumberGenerator is a mock of TxNumberGenerator interface.
type MockTxNumberGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockTxNumberGeneratorMockRecorder
}

// MockTxNumberGeneratorMockRecorder is the mock recorder for MockTxNumberGenerator.
type MockTxNumberGeneratorMockRecorder struct {
	mock *MockTxNumberGenerator
}

// NewMockTxNumberGenerator creates a new mock instance.
func NewMockTxNumberGenerator(ctrl *gomock.Controller) *MockTxNumberGenerator {
	mock := &MockTxNumberGenerator{ctrl: ctrl}
	mock.recorder = &MockTxNumberGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxNumberGenerator) EXPECT() *MockTxNumberGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockTxNumberGenerator) Generate() domain.TransactionNumber {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(domain.TransactionNumber)
	return ret0
}

// Generate indicates an expected call of Generate.
func (mr *MockTxNumberGeneratorMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockTxNumberGenerator)(nil).Generate))
}
