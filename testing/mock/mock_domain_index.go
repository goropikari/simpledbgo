// Code generated by MockGen. DO NOT EDIT.
// Source: index.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/goropikari/simpledbgo/domain"
)

// MockIndexer is a mock of Indexer interface.
type MockIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockIndexerMockRecorder
}

// MockIndexerMockRecorder is the mock recorder for MockIndexer.
type MockIndexerMockRecorder struct {
	mock *MockIndexer
}

// NewMockIndexer creates a new mock instance.
func NewMockIndexer(ctrl *gomock.Controller) *MockIndexer {
	mock := &MockIndexer{ctrl: ctrl}
	mock.recorder = &MockIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIndexer) EXPECT() *MockIndexerMockRecorder {
	return m.recorder
}

// BeforeFirst mocks base method.
func (m *MockIndexer) BeforeFirst(searchKey domain.Constant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeforeFirst", searchKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// BeforeFirst indicates an expected call of BeforeFirst.
func (mr *MockIndexerMockRecorder) BeforeFirst(searchKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeforeFirst", reflect.TypeOf((*MockIndexer)(nil).BeforeFirst), searchKey)
}

// Close mocks base method.
func (m *MockIndexer) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockIndexerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIndexer)(nil).Close))
}

// Delete mocks base method.
func (m *MockIndexer) Delete(arg0 domain.Constant, arg1 domain.RecordID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIndexerMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIndexer)(nil).Delete), arg0, arg1)
}

// GetDataRecordID mocks base method.
func (m *MockIndexer) GetDataRecordID() (domain.RecordID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDataRecordID")
	ret0, _ := ret[0].(domain.RecordID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDataRecordID indicates an expected call of GetDataRecordID.
func (mr *MockIndexerMockRecorder) GetDataRecordID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDataRecordID", reflect.TypeOf((*MockIndexer)(nil).GetDataRecordID))
}

// HasNext mocks base method.
func (m *MockIndexer) HasNext() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasNext")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasNext indicates an expected call of HasNext.
func (mr *MockIndexerMockRecorder) HasNext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasNext", reflect.TypeOf((*MockIndexer)(nil).HasNext))
}

// Insert mocks base method.
func (m *MockIndexer) Insert(arg0 domain.Constant, arg1 domain.RecordID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockIndexerMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockIndexer)(nil).Insert), arg0, arg1)
}