// Code generated by MockGen. DO NOT EDIT.
// Source: buffer_manager.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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

// FlushAll mocks base method.
func (m *MockBufferManager) FlushAll(txnum int) error {
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