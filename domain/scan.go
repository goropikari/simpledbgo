package domain

import (
	"errors"
	"fmt"
	"log"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// ErrNotUpdatable indicates scanner is not updatable.
var ErrNotUpdatable = errors.New("can't update query")

// Scanner is an interface of scanner.
type Scanner interface {
	MoveToFirst() error
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
	SetVal(FieldName, Constant) error
	SetInt32(FieldName, int32) error
	SetString(FieldName, string) error
	AdvanceNextInsertSlotID() error
	Delete() error
	RecordID() RecordID
	MoveToRecordID(rid RecordID) error
}

// Table is a model of database table.
type Table struct {
	txn           Transaction
	layout        *Layout
	recordPage    *RecordPage
	tblName       TableName
	currentSlotID SlotID
	err           error
}

// NewTable constructs a Table.
func NewTable(txn Transaction, tblName TableName, layout *Layout) (*Table, error) {
	tbl := &Table{
		txn:           txn,
		layout:        layout,
		tblName:       tblName,
		recordPage:    nil,
		currentSlotID: -1,
		err:           nil,
	}

	blkLen, err := txn.BlockLength(FileName(tblName))
	if err != nil {
		return nil, err
	}

	if blkLen == 0 {
		err := tbl.moveToNewBlock()
		if err != nil {
			return nil, err
		}
	} else {
		err := tbl.moveToBlock(0)
		if err != nil {
			return nil, err
		}
	}

	return tbl, nil
}

// Err returns iteration err.
func (tbl *Table) Err() error {
	return tbl.err
}

// Close closes the table.
func (tbl *Table) Close() {
	if tbl.recordPage != nil {
		tbl.txn.Unpin(tbl.recordPage.Block())
	}
}

// GetInt32 gets int32 from the table.
func (tbl *Table) GetInt32(fldName FieldName) (int32, error) {
	return tbl.recordPage.GetInt32(tbl.currentSlotID, fldName)
}

// GetString gets string from the table.
func (tbl *Table) GetString(fldName FieldName) (string, error) {
	return tbl.recordPage.GetString(tbl.currentSlotID, fldName)
}

// GetVal gets value from the table.
func (tbl *Table) GetVal(fldName FieldName) (Constant, error) {
	typ := tbl.layout.schema.Type(fldName)
	switch typ {
	case FInt32:
		val, err := tbl.GetInt32(fldName)
		if err != nil {
			return Constant{}, err
		}

		return NewConstant(VInt32, val), nil
	case FString:
		val, err := tbl.GetString(fldName)
		if err != nil {
			return Constant{}, err
		}

		return NewConstant(VString, val), nil
	case FUnknown:
		log.Fatal(errors.New("unexpected field type"))
	}

	return Constant{}, errors.New("GetVal error")
}

// SetInt32 sets int32 to the table.
func (tbl *Table) SetInt32(fldName FieldName, val int32) error {
	return tbl.recordPage.SetInt32(tbl.currentSlotID, fldName, val)
}

// SetString sets string to the table.
func (tbl *Table) SetString(fldName FieldName, val string) error {
	l := tbl.layout.Length(fldName)
	if len(val) > l {
		return fmt.Errorf("exceed varchar size %v: value '%v'", l, val)
	}

	return tbl.recordPage.SetString(tbl.currentSlotID, fldName, val)
}

// SetVal sets value to the table.
func (tbl *Table) SetVal(fldName FieldName, val Constant) error {
	typ := tbl.layout.schema.Type(fldName)
	switch typ {
	case FInt32:
		// TODO: check val type?
		err := tbl.SetInt32(fldName, val.ToInt32())
		if err != nil {
			return err
		}
	case FString:
		err := tbl.SetString(fldName, val.AsString())
		if err != nil {
			return err
		}
	case FUnknown:
		log.Fatal(errors.New("unexpected field type"))
	}

	return nil
}

// AdvanceNextInsertSlotID  advances current slot id to next to unused slot id.
// If there is no unused record, append file block.
func (tbl *Table) AdvanceNextInsertSlotID() error {
	slotID, err := tbl.recordPage.InsertAfter(tbl.currentSlotID)
	if err != nil {
		return err
	}
	tbl.currentSlotID = slotID

	for tbl.currentSlotID < 0 {
		last, err := tbl.isAtLastBlock()
		if err != nil {
			return err
		}
		if last {
			err = tbl.moveToNewBlock()
			if err != nil {
				return err
			}
		} else {
			blk := tbl.recordPage.Block()
			blkNum := blk.Number()
			err := tbl.moveToBlock(blkNum + 1)
			if err != nil {
				return err
			}
		}
		slotID, err := tbl.recordPage.InsertAfter(tbl.currentSlotID)
		if err != nil {
			return err
		}
		tbl.currentSlotID = slotID
	}

	return nil
}

// Delete deletes the current slot logically.
func (tbl *Table) Delete() error {
	return tbl.recordPage.Delete(tbl.currentSlotID)
}

// MoveToRecordID moves to the record id.
func (tbl *Table) MoveToRecordID(rid RecordID) error {
	tbl.Close()
	blk := NewBlock(FileName(tbl.tblName), rid.BlockNumber())
	recordPage, err := NewRecordPage(tbl.txn, blk, tbl.layout)
	if err != nil {
		return err
	}
	tbl.recordPage = recordPage
	tbl.currentSlotID = rid.SlotID()

	return nil
}

// RecordID is a identifier of record.
func (tbl *Table) RecordID() RecordID {
	blk := tbl.recordPage.Block()

	return NewRecordID(blk.Number(), tbl.currentSlotID)
}

// HasField checks the existence of the field.
func (tbl *Table) HasField(fldName FieldName) bool {
	return tbl.layout.schema.HasField(fldName)
}

// HasNext checks the existence of next record.
func (tbl *Table) HasNext() bool {
	currentSlotID, err := tbl.recordPage.NextAfter(tbl.currentSlotID)
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

		slotID, err := tbl.recordPage.NextAfter(tbl.currentSlotID)
		if err != nil {
			tbl.err = err

			return false
		}
		tbl.currentSlotID = slotID
	}

	return true
}

// isAtLastBlock checks whether the current block is last block or not.
func (tbl *Table) isAtLastBlock() (bool, error) {
	blk := tbl.recordPage.Block()
	size, err := tbl.txn.BlockLength(FileName(tbl.tblName))
	if err != nil {
		return false, err
	}

	blkNum, err := NewBlockNumber(size - 1)
	if err != nil {
		return false, err
	}

	return blk.Number() == blkNum, nil
}

func (tbl *Table) moveToNewBlock() error {
	tbl.Close()
	blk, err := tbl.txn.ExtendFile(FileName(tbl.tblName))
	if err != nil {
		return err
	}

	recordPage, err := NewRecordPage(tbl.txn, blk, tbl.layout)
	if err != nil {
		return err
	}

	tbl.recordPage = recordPage

	err = tbl.recordPage.Format()
	if err != nil {
		return err
	}

	tbl.currentSlotID = -1

	return nil
}

// MoveToFirst move to the first block of the table.
func (tbl *Table) MoveToFirst() error {
	return tbl.moveToBlock(0)
}

func (tbl *Table) moveToBlock(blkNum BlockNumber) error {
	tbl.Close()
	blk := NewBlock(FileName(tbl.tblName), blkNum)
	recordPage, err := NewRecordPage(tbl.txn, blk, tbl.layout)
	if err != nil {
		return err
	}

	tbl.recordPage = recordPage
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

	if err := scan.MoveToFirst(); err != nil {
		return nil, err
	}

	return scan, nil
}

// MoveToFirst moves to first record.
func (scan *ProductScan) MoveToFirst() error {
	err := scan.lhsScan.MoveToFirst()
	if err != nil {
		return err
	}

	scan.lhsScan.HasNext()
	if scan.lhsScan.Err() != nil {
		return scan.lhsScan.Err()
	}

	err = scan.rhsScan.MoveToFirst()
	if err != nil {
		return err
	}

	return nil
}

// HasNext checks the existence of next record.
func (plan *ProductScan) HasNext() bool {
	if plan.rhsScan.HasNext() {
		return true
	}
	if plan.rhsScan.Err() != nil {
		plan.err = plan.rhsScan.Err()

		return false
	}

	err := plan.rhsScan.MoveToFirst()
	if err != nil {
		plan.err = err

		return false
	}

	rhsFound := plan.rhsScan.HasNext()
	if plan.rhsScan.Err() != nil {
		plan.err = plan.rhsScan.Err()

		return false
	}

	lhsFound := plan.lhsScan.HasNext()
	if plan.lhsScan.Err() != nil {
		plan.err = plan.lhsScan.Err()

		return false
	}

	return lhsFound && rhsFound
}

// GetInt32 gets int32 from the table.
func (plan *ProductScan) GetInt32(fld FieldName) (int32, error) {
	if plan.lhsScan.HasField(fld) {
		return plan.lhsScan.GetInt32(fld)
	}

	if plan.rhsScan.HasField(fld) {
		return plan.rhsScan.GetInt32(fld)
	}

	return 0, fmt.Errorf("field %v not found", fld)
}

// GetString gets fld as string.
func (plan *ProductScan) GetString(fld FieldName) (string, error) {
	if plan.lhsScan.HasField(fld) {
		return plan.lhsScan.GetString(fld)
	}

	if plan.rhsScan.HasField(fld) {
		return plan.rhsScan.GetString(fld)
	}

	return "", fmt.Errorf("field %v not found", fld)
}

// GetVal gets fld as Constant.
func (plan *ProductScan) GetVal(fld FieldName) (Constant, error) {
	if plan.lhsScan.HasField(fld) {
		return plan.lhsScan.GetVal(fld)
	}

	if plan.rhsScan.HasField(fld) {
		return plan.rhsScan.GetVal(fld)
	}

	return Constant{}, fmt.Errorf("field %v not found", fld)
}

// HasField checks whether plan has fld as field or not.
func (plan *ProductScan) HasField(fld FieldName) bool {
	return plan.lhsScan.HasField(fld) || plan.rhsScan.HasField(fld)
}

// Close closes scan.
func (plan *ProductScan) Close() {
	plan.lhsScan.Close()
	plan.rhsScan.Close()
}

// Err returns error.
func (plan *ProductScan) Err() error {
	return plan.err
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

// MoveToFirst moves to first record.
func (s *SelectScan) MoveToFirst() error {
	return s.scan.MoveToFirst()
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

// GetInt32 gets int32 from the scanner.
func (s *SelectScan) GetInt32(fld FieldName) (int32, error) {
	return s.scan.GetInt32(fld)
}

// GetString gets string from the scanner.
func (s *SelectScan) GetString(fld FieldName) (string, error) {
	return s.scan.GetString(fld)
}

// GetVal gets value from the scanner.
func (s *SelectScan) GetVal(fld FieldName) (Constant, error) {
	return s.scan.GetVal(fld)
}

// HasField checks the existence of the field.
func (s *SelectScan) HasField(fld FieldName) bool {
	return s.scan.HasField(fld)
}

// Close closes scanner.
func (s *SelectScan) Close() {
	s.scan.Close()
}

// SetVal sets value at given fld.
func (s *SelectScan) SetVal(fld FieldName, c Constant) error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.SetVal(fld, c)
}

// SetInt32 sets int32 at given fld.
func (s *SelectScan) SetInt32(fld FieldName, x int32) error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.SetInt32(fld, x)
}

// SetString sets string at given fld.
func (s *SelectScan) SetString(fld FieldName, str string) error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.SetString(fld, str)
}

// AdvanceNextInsertSlotID returns insertion slot id.
func (s *SelectScan) AdvanceNextInsertSlotID() error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.AdvanceNextInsertSlotID()
}

// Delete deletes record.
func (s *SelectScan) Delete() error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.Delete()
}

// RecordID returns current record id.
func (s *SelectScan) RecordID() RecordID {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		panic(ErrNotUpdatable)
	}

	return us.RecordID()
}

// MoveToRecordID moves given RecordID.
func (s *SelectScan) MoveToRecordID(rid RecordID) error {
	us, ok := s.scan.(UpdateScanner)
	if !ok {
		return ErrNotUpdatable
	}

	return us.MoveToRecordID(rid)
}

// Err returns error.
func (s *SelectScan) Err() error {
	return s.err
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

// MoveToFirst moves to first record.
func (s *ProjectScan) MoveToFirst() error {
	return s.scan.MoveToFirst()
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

	return 0, fmt.Errorf("field %v not found", fld)
}

// GetString gets string from the scanner.
func (s *ProjectScan) GetString(fld FieldName) (string, error) {
	if s.HasField(fld) {
		return s.scan.GetString(fld)
	}

	return "", fmt.Errorf("field %v not found", fld)
}

// GetVal gets value from the scanner.
func (s *ProjectScan) GetVal(fld FieldName) (Constant, error) {
	if s.HasField(fld) {
		return s.scan.GetVal(fld)
	}

	return Constant{}, fmt.Errorf("field %v not found", fld)
}

// HasField checks the existence of the field.
func (s *ProjectScan) HasField(fld FieldName) bool {
	for _, f := range s.fields {
		if f == fld {
			return true
		}
	}

	return false
}

// Close closes scanner.
func (s *ProjectScan) Close() {
	s.scan.Close()
}

// Err returns error.
func (s *ProjectScan) Err() error {
	return s.err
}
