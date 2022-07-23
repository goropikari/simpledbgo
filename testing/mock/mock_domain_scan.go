// Code generated by MockGen. DO NOT EDIT.
// Source: scan.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/goropikari/simpledbgo/domain"
)

// MockScanner is a mock of Scanner interface.
type MockScanner struct {
	ctrl     *gomock.Controller
	recorder *MockScannerMockRecorder
}

// MockScannerMockRecorder is the mock recorder for MockScanner.
type MockScannerMockRecorder struct {
	mock *MockScanner
}

// NewMockScanner creates a new mock instance.
func NewMockScanner(ctrl *gomock.Controller) *MockScanner {
	mock := &MockScanner{ctrl: ctrl}
	mock.recorder = &MockScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScanner) EXPECT() *MockScannerMockRecorder {
	return m.recorder
}

// BeforeFirst mocks base method.
func (m *MockScanner) BeforeFirst() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeforeFirst")
	ret0, _ := ret[0].(error)
	return ret0
}

// BeforeFirst indicates an expected call of BeforeFirst.
func (mr *MockScannerMockRecorder) BeforeFirst() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeforeFirst", reflect.TypeOf((*MockScanner)(nil).BeforeFirst))
}

// Close mocks base method.
func (m *MockScanner) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockScannerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockScanner)(nil).Close))
}

// Err mocks base method.
func (m *MockScanner) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockScannerMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockScanner)(nil).Err))
}

// GetInt32 mocks base method.
func (m *MockScanner) GetInt32(arg0 domain.FieldName) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt32", arg0)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInt32 indicates an expected call of GetInt32.
func (mr *MockScannerMockRecorder) GetInt32(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt32", reflect.TypeOf((*MockScanner)(nil).GetInt32), arg0)
}

// GetString mocks base method.
func (m *MockScanner) GetString(arg0 domain.FieldName) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetString indicates an expected call of GetString.
func (mr *MockScannerMockRecorder) GetString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockScanner)(nil).GetString), arg0)
}

// GetVal mocks base method.
func (m *MockScanner) GetVal(arg0 domain.FieldName) (domain.Constant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVal", arg0)
	ret0, _ := ret[0].(domain.Constant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVal indicates an expected call of GetVal.
func (mr *MockScannerMockRecorder) GetVal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVal", reflect.TypeOf((*MockScanner)(nil).GetVal), arg0)
}

// HasField mocks base method.
func (m *MockScanner) HasField(arg0 domain.FieldName) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasField", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasField indicates an expected call of HasField.
func (mr *MockScannerMockRecorder) HasField(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasField", reflect.TypeOf((*MockScanner)(nil).HasField), arg0)
}

// HasNext mocks base method.
func (m *MockScanner) HasNext() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasNext")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasNext indicates an expected call of HasNext.
func (mr *MockScannerMockRecorder) HasNext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasNext", reflect.TypeOf((*MockScanner)(nil).HasNext))
}

// MockUpdateScanner is a mock of UpdateScanner interface.
type MockUpdateScanner struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateScannerMockRecorder
}

// MockUpdateScannerMockRecorder is the mock recorder for MockUpdateScanner.
type MockUpdateScannerMockRecorder struct {
	mock *MockUpdateScanner
}

// NewMockUpdateScanner creates a new mock instance.
func NewMockUpdateScanner(ctrl *gomock.Controller) *MockUpdateScanner {
	mock := &MockUpdateScanner{ctrl: ctrl}
	mock.recorder = &MockUpdateScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdateScanner) EXPECT() *MockUpdateScannerMockRecorder {
	return m.recorder
}

// AdvanceNextInsertSlotID mocks base method.
func (m *MockUpdateScanner) AdvanceNextInsertSlotID() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdvanceNextInsertSlotID")
	ret0, _ := ret[0].(error)
	return ret0
}

// AdvanceNextInsertSlotID indicates an expected call of AdvanceNextInsertSlotID.
func (mr *MockUpdateScannerMockRecorder) AdvanceNextInsertSlotID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdvanceNextInsertSlotID", reflect.TypeOf((*MockUpdateScanner)(nil).AdvanceNextInsertSlotID))
}

// BeforeFirst mocks base method.
func (m *MockUpdateScanner) BeforeFirst() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeforeFirst")
	ret0, _ := ret[0].(error)
	return ret0
}

// BeforeFirst indicates an expected call of BeforeFirst.
func (mr *MockUpdateScannerMockRecorder) BeforeFirst() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeforeFirst", reflect.TypeOf((*MockUpdateScanner)(nil).BeforeFirst))
}

// Close mocks base method.
func (m *MockUpdateScanner) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockUpdateScannerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockUpdateScanner)(nil).Close))
}

// Delete mocks base method.
func (m *MockUpdateScanner) Delete() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete")
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUpdateScannerMockRecorder) Delete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUpdateScanner)(nil).Delete))
}

// Err mocks base method.
func (m *MockUpdateScanner) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockUpdateScannerMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockUpdateScanner)(nil).Err))
}

// GetInt32 mocks base method.
func (m *MockUpdateScanner) GetInt32(arg0 domain.FieldName) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt32", arg0)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInt32 indicates an expected call of GetInt32.
func (mr *MockUpdateScannerMockRecorder) GetInt32(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt32", reflect.TypeOf((*MockUpdateScanner)(nil).GetInt32), arg0)
}

// GetString mocks base method.
func (m *MockUpdateScanner) GetString(arg0 domain.FieldName) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetString indicates an expected call of GetString.
func (mr *MockUpdateScannerMockRecorder) GetString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockUpdateScanner)(nil).GetString), arg0)
}

// GetVal mocks base method.
func (m *MockUpdateScanner) GetVal(arg0 domain.FieldName) (domain.Constant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVal", arg0)
	ret0, _ := ret[0].(domain.Constant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVal indicates an expected call of GetVal.
func (mr *MockUpdateScannerMockRecorder) GetVal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVal", reflect.TypeOf((*MockUpdateScanner)(nil).GetVal), arg0)
}

// HasField mocks base method.
func (m *MockUpdateScanner) HasField(arg0 domain.FieldName) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasField", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasField indicates an expected call of HasField.
func (mr *MockUpdateScannerMockRecorder) HasField(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasField", reflect.TypeOf((*MockUpdateScanner)(nil).HasField), arg0)
}

// HasNext mocks base method.
func (m *MockUpdateScanner) HasNext() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasNext")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasNext indicates an expected call of HasNext.
func (mr *MockUpdateScannerMockRecorder) HasNext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasNext", reflect.TypeOf((*MockUpdateScanner)(nil).HasNext))
}

// MoveToRecordID mocks base method.
func (m *MockUpdateScanner) MoveToRecordID(rid domain.RecordID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MoveToRecordID", rid)
	ret0, _ := ret[0].(error)
	return ret0
}

// MoveToRecordID indicates an expected call of MoveToRecordID.
func (mr *MockUpdateScannerMockRecorder) MoveToRecordID(rid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MoveToRecordID", reflect.TypeOf((*MockUpdateScanner)(nil).MoveToRecordID), rid)
}

// RecordID mocks base method.
func (m *MockUpdateScanner) RecordID() domain.RecordID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordID")
	ret0, _ := ret[0].(domain.RecordID)
	return ret0
}

// RecordID indicates an expected call of RecordID.
func (mr *MockUpdateScannerMockRecorder) RecordID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordID", reflect.TypeOf((*MockUpdateScanner)(nil).RecordID))
}

// SetInt32 mocks base method.
func (m *MockUpdateScanner) SetInt32(arg0 domain.FieldName, arg1 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetInt32", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInt32 indicates an expected call of SetInt32.
func (mr *MockUpdateScannerMockRecorder) SetInt32(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInt32", reflect.TypeOf((*MockUpdateScanner)(nil).SetInt32), arg0, arg1)
}

// SetString mocks base method.
func (m *MockUpdateScanner) SetString(arg0 domain.FieldName, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetString", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetString indicates an expected call of SetString.
func (mr *MockUpdateScannerMockRecorder) SetString(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetString", reflect.TypeOf((*MockUpdateScanner)(nil).SetString), arg0, arg1)
}

// SetVal mocks base method.
func (m *MockUpdateScanner) SetVal(arg0 domain.FieldName, arg1 domain.Constant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetVal", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetVal indicates an expected call of SetVal.
func (mr *MockUpdateScannerMockRecorder) SetVal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetVal", reflect.TypeOf((*MockUpdateScanner)(nil).SetVal), arg0, arg1)
}
