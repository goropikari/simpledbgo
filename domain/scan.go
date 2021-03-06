package domain

import (
	"fmt"

	"github.com/goropikari/simpledbgo/errors"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// ErrNotUpdatable indicates scanner is not updatable.
var ErrNotUpdatable = errors.New("can't update query")

// Scanner is an interface of scanner.
type Scanner interface {
	// BeforeFirst move to the position before the first record.
	BeforeFirst() error
	HasNext() bool
	GetInt32(FieldName) (int32, error)
	GetString(FieldName) (string, error)
	GetVal(FieldName) (Constant, error)
	HasField(FieldName) bool
	Close()
	Err() error
}

// UpdateScanner is an interface of updatable scanner.
type UpdateScanner interface {
	Scanner
	SetInt32(FieldName, int32) error
	SetString(FieldName, string) error
	SetVal(FieldName, Constant) error
	AdvanceNextInsertSlotID() error
	Delete() error
	RecordID() RecordID
	MoveToRecordID(rid RecordID) error
}

// TableScan is a model of database table.
type TableScan struct {
	txn           Transaction
	layout        *Layout
	recordPage    *RecordPage
	tblName       TableName
	currentSlotID SlotID
	err           error
}

// NewTableScan constructs a Table.
func NewTableScan(txn Transaction, tblName TableName, layout *Layout) (*TableScan, error) {
	tbl := &TableScan{
		txn:           txn,
		layout:        layout,
		tblName:       tblName,
		recordPage:    nil,
		currentSlotID: -1,
		err:           nil,
	}

	blkLen, err := txn.BlockLength(FileName(tblName))
	if err != nil {
		return nil, errors.Err(err, "BlockLength")
	}

	if blkLen == 0 {
		err := tbl.moveToNewBlock()
		if err != nil {
			return nil, errors.Wrap(err, "moveToNewBlock")
		}
	} else {
		err := tbl.BeforeFirst()
		if err != nil {
			return nil, errors.Wrap(err, "BeforeFirst")
		}
	}

	return tbl, nil
}

// BeforeFirst move to the position before the first record.
// BeforeFirst implements Scanner.
func (tbl *TableScan) BeforeFirst() error {
	return tbl.moveToBlock(0)
}

// HasNext checks the existence of next record.
// HasNext implements Scanner.
func (tbl *TableScan) HasNext() bool {
	currentSlotID, err := tbl.recordPage.NextUsedSlot(tbl.currentSlotID)
	if err != nil {
		tbl.err = err

		return false
	}
	tbl.currentSlotID = currentSlotID

	for tbl.currentSlotID < 0 {
		last, err := tbl.isAtLastBlock()
		if err != nil {
			tbl.err = err

			return false
		}
		if last {
			return false
		}

		blk := tbl.recordPage.Block()
		if err := tbl.moveToBlock(blk.Number() + 1); err != nil {
			tbl.err = err

			return false
		}

		slotID, err := tbl.recordPage.NextUsedSlot(tbl.currentSlotID)
		if err != nil {
			tbl.err = err

			return false
		}
		tbl.currentSlotID = slotID
	}

	return true
}

// GetInt32 gets int32 from the table.
// GetInt32  implements Scanner.
func (tbl *TableScan) GetInt32(fldName FieldName) (int32, error) {
	return tbl.recordPage.GetInt32(tbl.currentSlotID, fldName)
}

// GetString gets string from the table.
// GetString  implements Scanner.
func (tbl *TableScan) GetString(fldName FieldName) (string, error) {
	return tbl.recordPage.GetString(tbl.currentSlotID, fldName)
}

// GetVal gets value from the table.
// GetVal implements Scanner.
func (tbl *TableScan) GetVal(fldName FieldName) (Constant, error) {
	typ := tbl.layout.schema.Type(fldName)
	switch typ {
	case Int32FieldType:
		val, err := tbl.GetInt32(fldName)
		if err != nil {
			return Constant{}, errors.Err(err, "GetInt32")
		}

		return NewConstant(Int32FieldType, val), nil
	case StringFieldType:
		val, err := tbl.GetString(fldName)
		if err != nil {
			return Constant{}, errors.Err(err, "GetString")
		}

		return NewConstant(StringFieldType, val), nil
	case UnknownFieldType:
		return Constant{}, errors.New("unexpected field type")
	}

	return Constant{}, errors.New("GetVal error")
}

// HasField checks the existence of the field.
// HasField implements Scanner.
func (tbl *TableScan) HasField(fldName FieldName) bool {
	return tbl.layout.schema.HasField(fldName)
}

// Close closes the table.
// Close implements Scanner.
func (tbl *TableScan) Close() {
	if tbl.recordPage != nil {
		tbl.txn.Unpin(tbl.recordPage.Block())
	}
}

// Err returns iteration err.
// Err implements Scanner.
func (tbl *TableScan) Err() error {
	return tbl.err
}

// SetInt32 sets int32 to the table.
// SetInt32 implements UpdateScanner.
func (tbl *TableScan) SetInt32(fldName FieldName, val int32) error {
	return tbl.recordPage.SetInt32(tbl.currentSlotID, fldName, val)
}

// SetString sets string to the table.
// SetString implements UpdateScanner.
func (tbl *TableScan) SetString(fldName FieldName, val string) error {
	l := tbl.layout.Length(fldName)
	if len(val) > l {
		return fmt.Errorf("exceed varchar size %v: value '%v'", l, val)
	}

	return tbl.recordPage.SetString(tbl.currentSlotID, fldName, val)
}

// SetVal sets value to the table.
// SetVal implements UpdateScanner.
func (tbl *TableScan) SetVal(fldName FieldName, val Constant) error {
	typ := tbl.layout.schema.Type(fldName)
	switch typ {
	case Int32FieldType:
		v, err := val.AsInt32()
		if err != nil {
			return errors.Err(err, "AsInt32")
		}
		if err := tbl.SetInt32(fldName, v); err != nil {
			return errors.Err(err, "SetInt32")
		}
	case StringFieldType:
		v, err := val.AsString()
		if err != nil {
			return errors.Err(err, "AsString")
		}
		if err := tbl.SetString(fldName, v); err != nil {
			return errors.Err(err, "SetString")
		}
	case UnknownFieldType:
		return ErrUnsupportedFieldType
	}

	return nil
}

// AdvanceNextInsertSlotID  advances current slot id to next to unused slot id.
// If there is no unused record, append file block.
// AdvanceNextInsertSlotID implements UpdateScanner.
func (tbl *TableScan) AdvanceNextInsertSlotID() error {
	slotID, err := tbl.recordPage.InsertAfter(tbl.currentSlotID)
	if err != nil {
		return errors.Err(err, "InsertAfter")
	}
	tbl.currentSlotID = slotID

	for tbl.currentSlotID < 0 {
		last, err := tbl.isAtLastBlock()
		if err != nil {
			return errors.Err(err, "isAtLastBlock")
		}
		if last {
			err = tbl.moveToNewBlock()
			if err != nil {
				return errors.Err(err, "moveToNewBlock")
			}
		} else {
			blk := tbl.recordPage.Block()
			blkNum := blk.Number()
			err := tbl.moveToBlock(blkNum + 1)
			if err != nil {
				return errors.Err(err, "moveToBlock")
			}
		}
		slotID, err := tbl.recordPage.InsertAfter(tbl.currentSlotID)
		if err != nil {
			return errors.Err(err, "InsertAfter")
		}
		tbl.currentSlotID = slotID
	}

	return nil
}

// Delete deletes the current slot logically.
// Delete implements UpdateScanner.
func (tbl *TableScan) Delete() error {
	return tbl.recordPage.Delete(tbl.currentSlotID)
}

// RecordID is a identifier of record.
// RecordID implements UpdateScanner.
func (tbl *TableScan) RecordID() RecordID {
	blk := tbl.recordPage.Block()

	return NewRecordID(blk.Number(), tbl.currentSlotID)
}

// MoveToRecordID moves to the record id.
// MoveToRecordID implements UpdateScanner.
func (tbl *TableScan) MoveToRecordID(rid RecordID) error {
	tbl.Close()
	blk := NewBlock(FileName(tbl.tblName), rid.BlockNumber())
	recordPage, err := NewRecordPage(tbl.txn, blk, tbl.layout)
	if err != nil {
		return errors.Err(err, "NewRecordPage")
	}
	tbl.recordPage = recordPage
	tbl.currentSlotID = rid.SlotID()

	return nil
}

func (tbl *TableScan) moveToBlock(blkNum BlockNumber) error {
	tbl.Close()
	blk := NewBlock(FileName(tbl.tblName), blkNum)
	recordPage, err := NewRecordPage(tbl.txn, blk, tbl.layout)
	if err != nil {
		return errors.Err(err, "NewRecordPage")
	}

	tbl.recordPage = recordPage
	tbl.currentSlotID = -1

	return nil
}

// isAtLastBlock checks whether the current block is last block or not.
func (tbl *TableScan) isAtLastBlock() (bool, error) {
	blk := tbl.recordPage.Block()
	size, err := tbl.txn.BlockLength(FileName(tbl.tblName))
	if err != nil {
		return false, errors.Err(err, "txn.BlockLength")
	}

	blkNum, err := NewBlockNumber(size - 1)
	if err != nil {
		return false, errors.Err(err, "NewBlockNumber")
	}

	return blk.Number() == blkNum, nil
}

func (tbl *TableScan) moveToNewBlock() error {
	tbl.Close()
	blk, err := tbl.txn.ExtendFile(FileName(tbl.tblName))
	if err != nil {
		return errors.Err(err, "txn.ExtendFile")
	}

	recordPage, err := NewRecordPage(tbl.txn, blk, tbl.layout)
	if err != nil {
		return errors.Err(err, "NewRecordPage")
	}

	tbl.recordPage = recordPage

	err = tbl.recordPage.Format()
	if err != nil {
		return errors.Err(err, "Page.Format")
	}

	tbl.currentSlotID = -1

	return nil
}

// ProductScan is scanner of table product.
type ProductScan struct {
	lhsScan Scanner
	rhsScan Scanner
	err     error
}

// NewProductScan constructs a ProductScan.
func NewProductScan(lhs, rhs Scanner) (*ProductScan, error) {
	scan := &ProductScan{
		lhsScan: lhs,
		rhsScan: rhs,
	}

	if err := scan.BeforeFirst(); err != nil {
		return nil, errors.Err(err, "BeforeFirst")
	}

	return scan, nil
}

// BeforeFirst move to the position before the first record.
// BeforeFirst implements Scanner.
func (scan *ProductScan) BeforeFirst() error {
	err := scan.lhsScan.BeforeFirst()
	if err != nil {
		return errors.Err(err, "BeforeFirst")
	}

	scan.lhsScan.HasNext()
	if scan.lhsScan.Err() != nil {
		return errors.Wrap(scan.lhsScan.Err(), "failed to Scan")
	}

	err = scan.rhsScan.BeforeFirst()
	if err != nil {
		return errors.Err(err, "BeforeFirst")
	}

	return nil
}

// HasNext checks the existence of next record.
// HasNext implements Scanner.
func (scan *ProductScan) HasNext() bool {
	if scan.rhsScan.HasNext() {
		return true
	}
	if scan.rhsScan.Err() != nil {
		scan.err = scan.rhsScan.Err()

		return false
	}

	err := scan.rhsScan.BeforeFirst()
	if err != nil {
		scan.err = err

		return false
	}

	rhsFound := scan.rhsScan.HasNext()
	if scan.rhsScan.Err() != nil {
		scan.err = scan.rhsScan.Err()

		return false
	}

	lhsFound := scan.lhsScan.HasNext()
	if scan.lhsScan.Err() != nil {
		scan.err = scan.lhsScan.Err()

		return false
	}

	return lhsFound && rhsFound
}

// GetInt32 gets int32 from the table.
// GetInt32 implements Scanner.
func (scan *ProductScan) GetInt32(fld FieldName) (int32, error) {
	if scan.lhsScan.HasField(fld) {
		return scan.lhsScan.GetInt32(fld)
	}

	if scan.rhsScan.HasField(fld) {
		return scan.rhsScan.GetInt32(fld)
	}

	return 0, fieldNotFoudError(fld)
}

// GetString gets string from the table.
// GetString  implements Scanner.
func (scan *ProductScan) GetString(fld FieldName) (string, error) {
	if scan.lhsScan.HasField(fld) {
		return scan.lhsScan.GetString(fld)
	}

	if scan.rhsScan.HasField(fld) {
		return scan.rhsScan.GetString(fld)
	}

	return "", fieldNotFoudError(fld)
}

// GetVal gets value from the table.
// GetVal implements Scanner.
func (scan *ProductScan) GetVal(fld FieldName) (Constant, error) {
	if scan.lhsScan.HasField(fld) {
		return scan.lhsScan.GetVal(fld)
	}

	if scan.rhsScan.HasField(fld) {
		return scan.rhsScan.GetVal(fld)
	}

	return Constant{}, fieldNotFoudError(fld)
}

// HasField checks the existence of the field.
// HasField implements Scanner.
func (scan *ProductScan) HasField(fld FieldName) bool {
	return scan.lhsScan.HasField(fld) || scan.rhsScan.HasField(fld)
}

// Close closes the table.
// Close implements Scanner.
func (scan *ProductScan) Close() {
	scan.lhsScan.Close()
	scan.rhsScan.Close()
}

// Err returns iteration err.
// Err implements Scanner.
func (scan *ProductScan) Err() error {
	return scan.err
}

// SelectScan is scanner of select query.
type SelectScan struct {
	scan Scanner
	pred *Predicate
	err  error
}

// NewSelectScan constructs SelectScan.
func NewSelectScan(s Scanner, pred *Predicate) *SelectScan {
	return &SelectScan{
		scan: s,
		pred: pred,
	}
}

// BeforeFirst move to the position before the first record.
// BeforeFirst implements Scanner.
func (s *SelectScan) BeforeFirst() error {
	return s.scan.BeforeFirst()
}

// HasNext checks the existence of next record.
func (s *SelectScan) HasNext() bool {
	for s.scan.HasNext() {
		if s.pred.IsSatisfied(s.scan) {
			return true
		}
	}
	s.err = s.scan.Err()

	return false
}

// GetInt32 gets int32 from the table.
// GetInt32 implements Scanner.
func (s *SelectScan) GetInt32(fld FieldName) (int32, error) {
	return s.scan.GetInt32(fld)
}

// GetString gets string from the table.
// GetString  implements Scanner.
func (s *SelectScan) GetString(fld FieldName) (string, error) {
	return s.scan.GetString(fld)
}

// GetVal gets value from the table.
// GetVal implements Scanner.
func (s *SelectScan) GetVal(fld FieldName) (Constant, error) {
	return s.scan.GetVal(fld)
}

// HasField checks the existence of the field.
// HasField implements Scanner.
func (s *SelectScan) HasField(fld FieldName) bool {
	return s.scan.HasField(fld)
}

// Close closes the table.
// Close implements Scanner.
func (s *SelectScan) Close() {
	s.scan.Close()
}

// SetVal sets value to the table.
// SetVal implements UpdateScanner.
func (s *SelectScan) SetVal(fld FieldName, c Constant) error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.SetVal(fld, c)
}

// SetInt32 sets int32 to the table.
// SetInt32 implements UpdateScanner.
func (s *SelectScan) SetInt32(fld FieldName, x int32) error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.SetInt32(fld, x)
}

// SetString sets string to the table.
// SetString implements UpdateScanner.
func (s *SelectScan) SetString(fld FieldName, str string) error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.SetString(fld, str)
}

// AdvanceNextInsertSlotID  advances current slot id to next to unused slot id.
// If there is no unused record, append file block.
// AdvanceNextInsertSlotID implements UpdateScanner.
func (s *SelectScan) AdvanceNextInsertSlotID() error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.AdvanceNextInsertSlotID()
}

// Delete deletes the current slot logically.
// Delete implements UpdateScanner.
func (s *SelectScan) Delete() error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.Delete()
}

// RecordID is a identifier of record.
// RecordID implements UpdateScanner.
func (s *SelectScan) RecordID() RecordID {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		panic(ErrNotUpdatable)
	}

	return us.RecordID()
}

// MoveToRecordID moves to the record id.
// MoveToRecordID implements UpdateScanner.
func (s *SelectScan) MoveToRecordID(rid RecordID) error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.MoveToRecordID(rid)
}

// Err returns iteration err.
// Err implements Scanner.
func (s *SelectScan) Err() error {
	return s.err
}

// IndexSelectScan is scanner using index.
type IndexSelectScan struct {
	ts  *TableScan
	idx Indexer
	val Constant
	err error
}

// NewIndexSelectScan constructs a IndexSelectScan.
func NewIndexSelectScan(ts *TableScan, idx Indexer, val Constant) (*IndexSelectScan, error) {
	s := &IndexSelectScan{
		ts:  ts,
		idx: idx,
		val: val,
		err: nil,
	}
	if err := s.BeforeFirst(); err != nil {
		return nil, errors.Err(err, "BeforeFirst")
	}

	return s, nil
}

// BeforeFirst move to the position before the first record.
// BeforeFirst implements Scanner.
func (ss *IndexSelectScan) BeforeFirst() error {
	return ss.idx.BeforeFirst(ss.val)
}

// HasNext checks the existence of next record.
func (ss *IndexSelectScan) HasNext() bool {
	found := ss.idx.HasNext()
	if found {
		rid, err := ss.idx.GetDataRecordID()
		if err != nil {
			ss.err = errors.Err(err, "GetDataRecordID")

			return false
		}
		if err := ss.ts.MoveToRecordID(rid); err != nil {
			ss.err = errors.Err(err, "MoveToRecordID")

			return false
		}
	}
	if ss.idx.Err() != nil {
		ss.err = ss.idx.Err()
	}

	return found
}

// GetInt32 gets int32 from the table.
func (ss *IndexSelectScan) GetInt32(fldName FieldName) (int32, error) {
	return ss.ts.GetInt32(fldName)
}

// GetString gets string from the table.
// GetString  implements Scanner.
func (ss *IndexSelectScan) GetString(fldName FieldName) (string, error) {
	return ss.ts.GetString(fldName)
}

// GetVal gets value from the table.
// GetVal implements Scanner.
func (ss *IndexSelectScan) GetVal(fldName FieldName) (Constant, error) {
	return ss.ts.GetVal(fldName)
}

// HasField checks the existence of the field.
// HasField implements Scanner.
func (ss *IndexSelectScan) HasField(fldName FieldName) bool {
	return ss.ts.HasField(fldName)
}

// Close closes the table.
// Close implements Scanner.
func (ss *IndexSelectScan) Close() {
	ss.idx.Close()
	ss.ts.Close()
}

// Err returns iteration err.
// Err implements Scanner.
func (ss *IndexSelectScan) Err() error {
	return ss.err
}

// ProjectScan is scanner of projection scanner.
type ProjectScan struct {
	scan   Scanner
	fields []FieldName
	err    error
}

// NewProjectScan constructs a ProjectScan.
func NewProjectScan(s Scanner, fields []FieldName) *ProjectScan {
	return &ProjectScan{
		scan:   s,
		fields: fields,
	}
}

// BeforeFirst move to the position before the first record.
// BeforeFirst implements Scanner.
func (s *ProjectScan) BeforeFirst() error {
	return s.scan.BeforeFirst()
}

// HasNext checks the existence of next record.
func (s *ProjectScan) HasNext() bool {
	found := s.scan.HasNext()
	s.err = s.scan.Err()

	return found
}

// GetInt32 gets int32 from the table.
func (s *ProjectScan) GetInt32(fld FieldName) (int32, error) {
	if s.HasField(fld) {
		return s.scan.GetInt32(fld)
	}

	return 0, fieldNotFoudError(fld)
}

// GetString gets string from the table.
// GetString  implements Scanner.
func (s *ProjectScan) GetString(fld FieldName) (string, error) {
	if s.HasField(fld) {
		return s.scan.GetString(fld)
	}

	return "", fieldNotFoudError(fld)
}

// GetVal gets value from the table.
// GetVal implements Scanner.
func (s *ProjectScan) GetVal(fld FieldName) (Constant, error) {
	if s.HasField(fld) {
		return s.scan.GetVal(fld)
	}

	return Constant{}, fieldNotFoudError(fld)
}

// HasField checks the existence of the field.
// HasField implements Scanner.
func (s *ProjectScan) HasField(fld FieldName) bool {
	return s.scan.HasField(fld)
}

// Close closes the table.
// Close implements Scanner.
func (s *ProjectScan) Close() {
	s.scan.Close()
}

// Err returns error.
func (s *ProjectScan) Err() error {
	return s.err
}

func fieldNotFoudError(fld FieldName) error {
	return fmt.Errorf("column \"%v\" does not exist", fld)
}
