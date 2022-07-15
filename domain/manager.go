package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// FileManager is an interface of file manager.
type FileManager interface {
	CopyBlockToPage(Block, *Page) error
	CopyPageToBlock(*Page, Block) error
	BlockLength(FileName) (int32, error)
	ExtendFile(FileName) (Block, error)
	BlockSize() BlockSize
	CreatePage() (*Page, error)
	IsInit() bool
}

// LogManager is an interface of log manager.
type LogManager interface {
	FlushLSN(LSN) error
	Flush() error
	AppendRecord([]byte) (LSN, error)
	AppendNewBlock() (Block, error)
	Iterator() (LogIterator, error)
	LogFileName() FileName
}

// LogIterator is a iterator of log record.
type LogIterator interface {
	HasNext() bool
	Next() ([]byte, error)
}

// MetadataManager is an interface of MetadataManager.
type MetadataManager interface {
	CreateTable(tblName TableName, sch *Schema, txn Transaction) error
	GetTableLayout(tblName TableName, txn Transaction) (*Layout, error)
	CreateView(viewName ViewName, viewDef ViewDef, txn Transaction) error
	GetViewDef(viewName ViewName, txn Transaction) (ViewDef, error)
	CreateIndex(idxName IndexName, tblName TableName, fldName FieldName, txn Transaction) error
	GetIndexInfo(tblName TableName, txn Transaction) (map[FieldName]*IndexInfo, error)
	GetStatInfo(tblName TableName, layout *Layout, txn Transaction) (StatInfo, error)
}

// BufferManager is an interface of buffer manager.
type BufferPoolManager interface {
	Available() int
	FlushAll(txnum TransactionNumber) error
	Unpin(buf *Buffer)
	Pin(Block) (*Buffer, error)
}

// ConcurrencyManager is an interface of concurrency manager.
type ConcurrencyManager interface {
	SLock(Block) error
	XLock(Block) error
	Release()
}
