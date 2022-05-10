package domain

const (
	fldBlock   = "block"
	fldID      = "id"
	fldDataVal = "dataval"
)

// IndexInfo is a model of information of index.
type IndexInfo struct {
	gen       IndexGenerator
	cal       SearchCostCalculator
	idxName   IndexName
	fldName   FieldName
	txn       Transaction
	tblSchema *Schema
	layout    *Layout
	statInfo  StatInfo
}

// NewIndexInfo constructs an IndexInfo.
func NewIndexInfo(factory IndexFactory, idxName IndexName, fldName FieldName, tblSchema *Schema, txn Transaction, si StatInfo) *IndexInfo {
	gen, cal := factory.Create()

	return &IndexInfo{
		gen:       gen,
		cal:       cal,
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
	return info.gen.Create(info.txn, info.idxName, info.layout)
}

// EstBlockAccessed estimates the number of accessing blocks.
func (info *IndexInfo) EstBlockAccessed() int {
	rpb := int(info.txn.BlockSize()) / int(info.layout.SlotSize())
	numBlks := info.statInfo.EstNumRecord() / rpb

	return info.cal.Calculate(numBlks, rpb)
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
	sch.AddInt32Field(fldBlock)
	sch.AddInt32Field(fldID)
	if tblSchema.Type(fldName) == FInt32 {
		sch.AddInt32Field(fldDataVal)
	} else {
		fldLen := tblSchema.Length(fldName)
		sch.AddStringField(fldDataVal, fldLen)
	}

	return NewLayout(sch)
}
