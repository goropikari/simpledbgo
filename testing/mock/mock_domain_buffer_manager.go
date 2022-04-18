// Code generated by MockGen. DO NOT EDIT.
// Source: buffer_manager.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/goropikari/simpledbgo/backend/domain"
)

// MockBufferManager is a mock of BufferManager interface.
type MockBufferManager struct {
	ctrl     *gomock.Controller
	recorder *MockBufferManagerMockRecorder
}

// MockBufferManagerMockRecorder is the mock recorder for MockBufferManager.
type MockBufferManagerMockRecorder struct {
	mock *MockBufferManager
}

// NewMockBufferManager creates a new mock instance.
func NewMockBufferManager(ctrl *gomock.Controller) *MockBufferManager {
	mock := &MockBufferManager{ctrl: ctrl}
	mock.recorder = &MockBufferManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBufferManager) EXPECT() *MockBufferManagerMockRecorder {
	return m.recorder
}

// Available mocks base method.
func (m *MockBufferManager) Available() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Available")
	ret0, _ := ret[0].(int)
	return ret0
}

// Available indicates an expected call of Available.
func (mr *MockBufferManagerMockRecorder) Available() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Available", reflect.TypeOf((*MockBufferManager)(nil).Available))
}

// FlushAll mocks base method.
func (m *MockBufferManager) FlushAll(txnum domain.TransactionNumber) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlushAll", txnum)
	ret0, _ := ret[0].(error)
	return ret0
}

// FlushAll indicates an expected call of FlushAll.
func (mr *MockBufferManagerMockRecorder) FlushAll(txnum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushAll", reflect.TypeOf((*MockBufferManager)(nil).FlushAll), txnum)
}

// Pin mocks base method.
func (m *MockBufferManager) Pin(arg0 domain.Block) (*domain.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pin", arg0)
	ret0, _ := ret[0].(*domain.Buffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Pin indicates an expected call of Pin.
func (mr *MockBufferManagerMockRecorder) Pin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pin", reflect.TypeOf((*MockBufferManager)(nil).Pin), arg0)
}

// Unpin mocks base method.
func (m *MockBufferManager) Unpin(buf *domain.Buffer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unpin", buf)
}

// Unpin indicates an expected call of Unpin.
func (mr *MockBufferManagerMockRecorder) Unpin(buf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unpin", reflect.TypeOf((*MockBufferManager)(nil).Unpin), buf)
}
