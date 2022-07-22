package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

const (
	// FldBlock is column name for block.
	FldBlock = "block"

	// FldID is column name for column.
	FldID = "id"

	// FldDataVal is column name for data value.
	FldDataVal = "dataval"
)

// SearchCostCalculator calculate search cost.
type SearchCostCalculator interface {
	Calculate(numBlk int, rpb int) int
}

// Indexer is an interface of index.
type Indexer interface {
	BeforeFirst(searchKey Constant) error
	HasNext() bool
	GetDataRecordID() (RecordID, error)
	Insert(Constant, RecordID) error
	Delete(Constant, RecordID) error
	Close()
	Err() error
}

// IndexDriver is driver for index.
type IndexDriver struct {
	fty IndexFactory
	cal SearchCostCalculator
}

// NewIndexDriver constructs a IndexDriver.
func NewIndexDriver(fty IndexFactory, cal SearchCostCalculator) IndexDriver {
	return IndexDriver{
		fty: fty,
		cal: cal,
	}
}

// Create creates a index.
func (d IndexDriver) Create(txn Transaction, name IndexName, layout *Layout) Indexer {
	return d.fty.Create(txn, name, layout)
}

// Calculate calculates a search cost.
func (d IndexDriver) Calculate(numBlk, rpb int) int {
	return d.cal.Calculate(numBlk, rpb)
}

// IndexFactory generates Index.
type IndexFactory interface {
	Create(Transaction, IndexName, *Layout) Indexer
}

// IndexName is a value object of index name.
type IndexName string

// NewIndexName constructs IndexName.
func NewIndexName(name string) (IndexName, error) {
	if len(name) > MaxFieldNameLength {
		return "", ErrExceedMaxFieldNameLength
	}

	return IndexName(name), nil
}

// String stringfies name.
func (name IndexName) String() string {
	return string(name)
}

// IndexInfo is a model of information of index.
type IndexInfo struct {
	driver    IndexDriver
	idxName   IndexName
	fldName   FieldName
	txn       Transaction
	tblSchema *Schema
	layout    *Layout
	statInfo  StatInfo
}

// NewIndexInfo constructs an IndexInfo.
func NewIndexInfo(driver IndexDriver, idxName IndexName, fldName FieldName, tblSchema *Schema, txn Transaction, si StatInfo) *IndexInfo {
	return &IndexInfo{
		driver:    driver,
		idxName:   idxName,
		fldName:   fldName,
		txn:       txn,
		tblSchema: tblSchema,
		layout:    createIdxLayout(tblSchema, fldName),
		statInfo:  si,
	}
}

// Open opens the index.
func (info *IndexInfo) Open() Indexer {
	return info.driver.Create(info.txn, info.idxName, info.layout)
}

// EstBlockAccessed estimates the number of accessing blocks.
func (info *IndexInfo) EstBlockAccessed() int {
	rpb := int(info.txn.BlockSize()) / int(info.layout.SlotSize())
	numBlks := info.statInfo.EstNumRecord() / rpb

	return info.driver.Calculate(numBlks, rpb)
}

// EstNumRecord estimates the number of records.
func (info *IndexInfo) EstNumRecord() int {
	return info.statInfo.EstNumRecord() / info.statInfo.EstDistinctVals(info.fldName)
}

// EstDistinctVals returns the estimation of the number of distinct values.
func (info *IndexInfo) EstDistinctVals(fldName FieldName) int {
	if info.fldName == fldName {
		return 1
	}

	return info.statInfo.EstDistinctVals(fldName)
}

func createIdxLayout(tblSchema *Schema, fldName FieldName) *Layout {
	sch := NewSchema()
	sch.AddInt32Field(FldBlock)
	sch.AddInt32Field(FldID)

	switch tblSchema.Type(fldName) {
	case Int32FieldType:
		sch.AddInt32Field(FldDataVal)
	case StringFieldType:
		fldLen := tblSchema.Length(fldName)
		sch.AddStringField(FldDataVal, fldLen)
	case UnknownFieldType:
		panic(ErrUnsupportedFieldType)
	default:
		panic(ErrUnsupportedFieldType)
	}

	return NewLayout(sch)
}
