// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/goropikari/simpledbgo/domain"
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

// BlockLength mocks base method.
func (m *MockFileManager) BlockLength(arg0 domain.FileName) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockLength", arg0)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockLength indicates an expected call of BlockLength.
func (mr *MockFileManagerMockRecorder) BlockLength(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockLength", reflect.TypeOf((*MockFileManager)(nil).BlockLength), arg0)
}

// BlockSize mocks base method.
func (m *MockFileManager) BlockSize() domain.BlockSize {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockSize")
	ret0, _ := ret[0].(domain.BlockSize)
	return ret0
}

// BlockSize indicates an expected call of BlockSize.
func (mr *MockFileManagerMockRecorder) BlockSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockSize", reflect.TypeOf((*MockFileManager)(nil).BlockSize))
}

// CopyBlockToPage mocks base method.
func (m *MockFileManager) CopyBlockToPage(arg0 domain.Block, arg1 *domain.Page) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyBlockToPage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyBlockToPage indicates an expected call of CopyBlockToPage.
func (mr *MockFileManagerMockRecorder) CopyBlockToPage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyBlockToPage", reflect.TypeOf((*MockFileManager)(nil).CopyBlockToPage), arg0, arg1)
}

// CopyPageToBlock mocks base method.
func (m *MockFileManager) CopyPageToBlock(arg0 *domain.Page, arg1 domain.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyPageToBlock", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyPageToBlock indicates an expected call of CopyPageToBlock.
func (mr *MockFileManagerMockRecorder) CopyPageToBlock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyPageToBlock", reflect.TypeOf((*MockFileManager)(nil).CopyPageToBlock), arg0, arg1)
}

// CreatePage mocks base method.
func (m *MockFileManager) CreatePage() (*domain.Page, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePage")
	ret0, _ := ret[0].(*domain.Page)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePage indicates an expected call of CreatePage.
func (mr *MockFileManagerMockRecorder) CreatePage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePage", reflect.TypeOf((*MockFileManager)(nil).CreatePage))
}

// ExtendFile mocks base method.
func (m *MockFileManager) ExtendFile(arg0 domain.FileName) (domain.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExtendFile", arg0)
	ret0, _ := ret[0].(domain.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExtendFile indicates an expected call of ExtendFile.
func (mr *MockFileManagerMockRecorder) ExtendFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtendFile", reflect.TypeOf((*MockFileManager)(nil).ExtendFile), arg0)
}

// IsInit mocks base method.
func (m *MockFileManager) IsInit() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsInit")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsInit indicates an expected call of IsInit.
func (mr *MockFileManagerMockRecorder) IsInit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsInit", reflect.TypeOf((*MockFileManager)(nil).IsInit))
}

// MockLogManager is a mock of LogManager interface.
type MockLogManager struct {
	ctrl     *gomock.Controller
	recorder *MockLogManagerMockRecorder
}

// MockLogManagerMockRecorder is the mock recorder for MockLogManager.
type MockLogManagerMockRecorder struct {
	mock *MockLogManager
}

// NewMockLogManager creates a new mock instance.
func NewMockLogManager(ctrl *gomock.Controller) *MockLogManager {
	mock := &MockLogManager{ctrl: ctrl}
	mock.recorder = &MockLogManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogManager) EXPECT() *MockLogManagerMockRecorder {
	return m.recorder
}

// AppendNewBlock mocks base method.
func (m *MockLogManager) AppendNewBlock() (domain.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendNewBlock")
	ret0, _ := ret[0].(domain.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AppendNewBlock indicates an expected call of AppendNewBlock.
func (mr *MockLogManagerMockRecorder) AppendNewBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendNewBlock", reflect.TypeOf((*MockLogManager)(nil).AppendNewBlock))
}

// AppendRecord mocks base method.
func (m *MockLogManager) AppendRecord(arg0 []byte) (domain.LSN, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendRecord", arg0)
	ret0, _ := ret[0].(domain.LSN)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AppendRecord indicates an expected call of AppendRecord.
func (mr *MockLogManagerMockRecorder) AppendRecord(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendRecord", reflect.TypeOf((*MockLogManager)(nil).AppendRecord), arg0)
}

// Flush mocks base method.
func (m *MockLogManager) Flush() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flush")
	ret0, _ := ret[0].(error)
	return ret0
}

// Flush indicates an expected call of Flush.
func (mr *MockLogManagerMockRecorder) Flush() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockLogManager)(nil).Flush))
}

// FlushLSN mocks base method.
func (m *MockLogManager) FlushLSN(arg0 domain.LSN) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlushLSN", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// FlushLSN indicates an expected call of FlushLSN.
func (mr *MockLogManagerMockRecorder) FlushLSN(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushLSN", reflect.TypeOf((*MockLogManager)(nil).FlushLSN), arg0)
}

// Iterator mocks base method.
func (m *MockLogManager) Iterator() (domain.LogIterator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Iterator")
	ret0, _ := ret[0].(domain.LogIterator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Iterator indicates an expected call of Iterator.
func (mr *MockLogManagerMockRecorder) Iterator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Iterator", reflect.TypeOf((*MockLogManager)(nil).Iterator))
}

// LogFileName mocks base method.
func (m *MockLogManager) LogFileName() domain.FileName {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogFileName")
	ret0, _ := ret[0].(domain.FileName)
	return ret0
}

// LogFileName indicates an expected call of LogFileName.
func (mr *MockLogManagerMockRecorder) LogFileName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogFileName", reflect.TypeOf((*MockLogManager)(nil).LogFileName))
}

// MockLogIterator is a mock of LogIterator interface.
type MockLogIterator struct {
	ctrl     *gomock.Controller
	recorder *MockLogIteratorMockRecorder
}

// MockLogIteratorMockRecorder is the mock recorder for MockLogIterator.
type MockLogIteratorMockRecorder struct {
	mock *MockLogIterator
}

// NewMockLogIterator creates a new mock instance.
func NewMockLogIterator(ctrl *gomock.Controller) *MockLogIterator {
	mock := &MockLogIterator{ctrl: ctrl}
	mock.recorder = &MockLogIteratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogIterator) EXPECT() *MockLogIteratorMockRecorder {
	return m.recorder
}

// Err mocks base method.
func (m *MockLogIterator) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockLogIteratorMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockLogIterator)(nil).Err))
}

// HasNext mocks base method.
func (m *MockLogIterator) HasNext() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasNext")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasNext indicates an expected call of HasNext.
func (mr *MockLogIteratorMockRecorder) HasNext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasNext", reflect.TypeOf((*MockLogIterator)(nil).HasNext))
}

// Next mocks base method.
func (m *MockLogIterator) Next() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Next indicates an expected call of Next.
func (mr *MockLogIteratorMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockLogIterator)(nil).Next))
}

// MockMetadataManager is a mock of MetadataManager interface.
type MockMetadataManager struct {
	ctrl     *gomock.Controller
	recorder *MockMetadataManagerMockRecorder
}

// MockMetadataManagerMockRecorder is the mock recorder for MockMetadataManager.
type MockMetadataManagerMockRecorder struct {
	mock *MockMetadataManager
}

// NewMockMetadataManager creates a new mock instance.
func NewMockMetadataManager(ctrl *gomock.Controller) *MockMetadataManager {
	mock := &MockMetadataManager{ctrl: ctrl}
	mock.recorder = &MockMetadataManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetadataManager) EXPECT() *MockMetadataManagerMockRecorder {
	return m.recorder
}

// CreateIndex mocks base method.
func (m *MockMetadataManager) CreateIndex(idxName domain.IndexName, tblName domain.TableName, fldName domain.FieldName, txn domain.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndex", idxName, tblName, fldName, txn)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIndex indicates an expected call of CreateIndex.
func (mr *MockMetadataManagerMockRecorder) CreateIndex(idxName, tblName, fldName, txn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndex", reflect.TypeOf((*MockMetadataManager)(nil).CreateIndex), idxName, tblName, fldName, txn)
}

// CreateTable mocks base method.
func (m *MockMetadataManager) CreateTable(tblName domain.TableName, sch *domain.Schema, txn domain.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTable", tblName, sch, txn)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTable indicates an expected call of CreateTable.
func (mr *MockMetadataManagerMockRecorder) CreateTable(tblName, sch, txn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTable", reflect.TypeOf((*MockMetadataManager)(nil).CreateTable), tblName, sch, txn)
}

// CreateView mocks base method.
func (m *MockMetadataManager) CreateView(viewName domain.ViewName, viewDef domain.ViewDef, txn domain.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateView", viewName, viewDef, txn)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateView indicates an expected call of CreateView.
func (mr *MockMetadataManagerMockRecorder) CreateView(viewName, viewDef, txn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateView", reflect.TypeOf((*MockMetadataManager)(nil).CreateView), viewName, viewDef, txn)
}

// GetIndexInfo mocks base method.
func (m *MockMetadataManager) GetIndexInfo(tblName domain.TableName, txn domain.Transaction) (map[domain.FieldName]*domain.IndexInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIndexInfo", tblName, txn)
	ret0, _ := ret[0].(map[domain.FieldName]*domain.IndexInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIndexInfo indicates an expected call of GetIndexInfo.
func (mr *MockMetadataManagerMockRecorder) GetIndexInfo(tblName, txn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIndexInfo", reflect.TypeOf((*MockMetadataManager)(nil).GetIndexInfo), tblName, txn)
}

// GetStatInfo mocks base method.
func (m *MockMetadataManager) GetStatInfo(tblName domain.TableName, layout *domain.Layout, txn domain.Transaction) (domain.StatInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatInfo", tblName, layout, txn)
	ret0, _ := ret[0].(domain.StatInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatInfo indicates an expected call of GetStatInfo.
func (mr *MockMetadataManagerMockRecorder) GetStatInfo(tblName, layout, txn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatInfo", reflect.TypeOf((*MockMetadataManager)(nil).GetStatInfo), tblName, layout, txn)
}

// GetTableLayout mocks base method.
func (m *MockMetadataManager) GetTableLayout(tblName domain.TableName, txn domain.Transaction) (*domain.Layout, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTableLayout", tblName, txn)
	ret0, _ := ret[0].(*domain.Layout)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTableLayout indicates an expected call of GetTableLayout.
func (mr *MockMetadataManagerMockRecorder) GetTableLayout(tblName, txn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTableLayout", reflect.TypeOf((*MockMetadataManager)(nil).GetTableLayout), tblName, txn)
}

// GetViewDef mocks base method.
func (m *MockMetadataManager) GetViewDef(viewName domain.ViewName, txn domain.Transaction) (domain.ViewDef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetViewDef", viewName, txn)
	ret0, _ := ret[0].(domain.ViewDef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetViewDef indicates an expected call of GetViewDef.
func (mr *MockMetadataManagerMockRecorder) GetViewDef(viewName, txn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetViewDef", reflect.TypeOf((*MockMetadataManager)(nil).GetViewDef), viewName, txn)
}

// MockBufferPoolManager is a mock of BufferPoolManager interface.
type MockBufferPoolManager struct {
	ctrl     *gomock.Controller
	recorder *MockBufferPoolManagerMockRecorder
}

// MockBufferPoolManagerMockRecorder is the mock recorder for MockBufferPoolManager.
type MockBufferPoolManagerMockRecorder struct {
	mock *MockBufferPoolManager
}

// NewMockBufferPoolManager creates a new mock instance.
func NewMockBufferPoolManager(ctrl *gomock.Controller) *MockBufferPoolManager {
	mock := &MockBufferPoolManager{ctrl: ctrl}
	mock.recorder = &MockBufferPoolManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBufferPoolManager) EXPECT() *MockBufferPoolManagerMockRecorder {
	return m.recorder
}

// Available mocks base method.
func (m *MockBufferPoolManager) Available() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Available")
	ret0, _ := ret[0].(int)
	return ret0
}

// Available indicates an expected call of Available.
func (mr *MockBufferPoolManagerMockRecorder) Available() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Available", reflect.TypeOf((*MockBufferPoolManager)(nil).Available))
}

// FlushAll mocks base method.
func (m *MockBufferPoolManager) FlushAll(txnum domain.TransactionNumber) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlushAll", txnum)
	ret0, _ := ret[0].(error)
	return ret0
}

// FlushAll indicates an expected call of FlushAll.
func (mr *MockBufferPoolManagerMockRecorder) FlushAll(txnum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushAll", reflect.TypeOf((*MockBufferPoolManager)(nil).FlushAll), txnum)
}

// Pin mocks base method.
func (m *MockBufferPoolManager) Pin(arg0 domain.Block) (*domain.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pin", arg0)
	ret0, _ := ret[0].(*domain.Buffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Pin indicates an expected call of Pin.
func (mr *MockBufferPoolManagerMockRecorder) Pin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pin", reflect.TypeOf((*MockBufferPoolManager)(nil).Pin), arg0)
}

// Unpin mocks base method.
func (m *MockBufferPoolManager) Unpin(buf *domain.Buffer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unpin", buf)
}

// Unpin indicates an expected call of Unpin.
func (mr *MockBufferPoolManagerMockRecorder) Unpin(buf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unpin", reflect.TypeOf((*MockBufferPoolManager)(nil).Unpin), buf)
}
