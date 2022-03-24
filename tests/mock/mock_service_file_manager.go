// Code generated by MockGen. DO NOT EDIT.
// Source: file_manager.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	core "github.com/goropikari/simpledb_go/backend/core"
)

// MockFileManager is a mock of FileManager interface.
type MockFileManager struct {
	ctrl     *gomock.Controller
	recorder *MockFileManagerMockRecorder
}

// MockFileManagerMockRecorder is the mock recorder for MockFileManager.
type MockFileManagerMockRecorder struct {
	mock *MockFileManager
}

// NewMockFileManager creates a new mock instance.
func NewMockFileManager(ctrl *gomock.Controller) *MockFileManager {
	mock := &MockFileManager{ctrl: ctrl}
	mock.recorder = &MockFileManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileManager) EXPECT() *MockFileManagerMockRecorder {
	return m.recorder
}

// AppendBlock mocks base method.
func (m *MockFileManager) AppendBlock(filename core.FileName) (*core.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendBlock", filename)
	ret0, _ := ret[0].(*core.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AppendBlock indicates an expected call of AppendBlock.
func (mr *MockFileManagerMockRecorder) AppendBlock(filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendBlock", reflect.TypeOf((*MockFileManager)(nil).AppendBlock), filename)
}

// CopyBlockToPage mocks base method.
func (m *MockFileManager) CopyBlockToPage(block *core.Block, page *core.Page) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyBlockToPage", block, page)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyBlockToPage indicates an expected call of CopyBlockToPage.
func (mr *MockFileManagerMockRecorder) CopyBlockToPage(block, page interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyBlockToPage", reflect.TypeOf((*MockFileManager)(nil).CopyBlockToPage), block, page)
}

// CopyPageToBlock mocks base method.
func (m *MockFileManager) CopyPageToBlock(page *core.Page, block *core.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyPageToBlock", page, block)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyPageToBlock indicates an expected call of CopyPageToBlock.
func (mr *MockFileManagerMockRecorder) CopyPageToBlock(page, block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyPageToBlock", reflect.TypeOf((*MockFileManager)(nil).CopyPageToBlock), page, block)
}

// GetBlockSize mocks base method.
func (m *MockFileManager) GetBlockSize() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetBlockSize indicates an expected call of GetBlockSize.
func (mr *MockFileManagerMockRecorder) GetBlockSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockSize", reflect.TypeOf((*MockFileManager)(nil).GetBlockSize))
}

// IsZero mocks base method.
func (m *MockFileManager) IsZero() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsZero")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsZero indicates an expected call of IsZero.
func (mr *MockFileManagerMockRecorder) IsZero() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsZero", reflect.TypeOf((*MockFileManager)(nil).IsZero))
}

// LastBlock mocks base method.
func (m *MockFileManager) LastBlock(filename core.FileName) (*core.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastBlock", filename)
	ret0, _ := ret[0].(*core.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastBlock indicates an expected call of LastBlock.
func (mr *MockFileManagerMockRecorder) LastBlock(filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastBlock", reflect.TypeOf((*MockFileManager)(nil).LastBlock), filename)
}

// PreparePage mocks base method.
func (m *MockFileManager) PreparePage() (*core.Page, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreparePage")
	ret0, _ := ret[0].(*core.Page)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PreparePage indicates an expected call of PreparePage.
func (mr *MockFileManagerMockRecorder) PreparePage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreparePage", reflect.TypeOf((*MockFileManager)(nil).PreparePage))
}
