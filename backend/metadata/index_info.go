package metadata

import (
	"github.com/goropikari/simpledbgo/backend/index"
	"github.com/goropikari/simpledbgo/backend/index/hash"
	"github.com/goropikari/simpledbgo/domain"
)

const (
	fldBlock   = "block"
	fldID      = "id"
	fldDataVal = "dataval"
)

// IndexInfo is a model of information of index.
type IndexInfo struct {
	idxName   domain.IndexName
	fldName   domain.FieldName
	txn       domain.Transaction
	tblSchema *domain.Schema
	layout    *domain.Layout
	statInfo  StatInfo
}

// NewIndexInfo constructs an IndexInfo.
func NewIndexInfo(idxName domain.IndexName, fldName domain.FieldName, tblSchema *domain.Schema, txn domain.Transaction, si StatInfo) *IndexInfo {
	return &IndexInfo{
		idxName:   idxName,
		fldName:   fldName,
		txn:       txn,
		tblSchema: tblSchema,
		layout:    createIdxLayout(tblSchema, fldName),
		statInfo:  si,
	}
}

// Open opens the index.
func (info *IndexInfo) Open() index.Index {
	return hash.NewIndex(info.txn, info.idxName, info.layout)
}

// EstBlockAccessed estimates the number of accessing blocks.
func (info *IndexInfo) EstBlockAccessed() int {
	rpb := int(info.txn.BlockSize()) / int(info.layout.SlotSize())
	numBlks := info.statInfo.EstNumRecord() / rpb

	return hash.SearchCost(numBlks, rpb)
}

// EstNumRecord estimates the number of records.
func (info *IndexInfo) EstNumRecord() int {
	return info.statInfo.EstNumRecord() / info.statInfo.EstDistinctVals(info.fldName)
}

// EstDistinctVals returns the estimation of the number of distinct values.
func (info *IndexInfo) EstDistinctVals(fldName domain.FieldName) int {
	if info.fldName == fldName {
		return 1
	}

	return info.statInfo.EstDistinctVals(fldName)
}

func createIdxLayout(tblSchema *domain.Schema, fldName domain.FieldName) *domain.Layout {
	sch := domain.NewSchema()
	sch.AddInt32Field(fldBlock)
	sch.AddInt32Field(fldID)
	if tblSchema.Type(fldName) == domain.FInt32 {
		sch.AddInt32Field(fldDataVal)
	} else {
		fldLen := tblSchema.Length(fldName)
		sch.AddStringField(fldDataVal, fldLen)
	}

	return domain.NewLayout(sch)
}
