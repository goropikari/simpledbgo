package domain

import (
	"log"

	"github.com/goropikari/simpledbgo/common"
	"github.com/goropikari/simpledbgo/errors"
)

// ErrFieldNotFound is an error that means specified field is not found.
var ErrFieldNotFound = errors.New("specified field is not found")

const (
	// RecordOffset is offset of record.
	RecordOffset = common.Int32Length
)

// Schema is model of table schema.
type Schema struct {
	fields []FieldName
	info   map[FieldName]*FieldInfo
}

// NewSchema constructs a Schema.
func NewSchema() *Schema {
	return &Schema{
		fields: make([]FieldName, 0),
		info:   make(map[FieldName]*FieldInfo),
	}
}

// HasField checks existence of fldname.
func (schema *Schema) HasField(fldname FieldName) bool {
	_, found := schema.info[fldname]

	return found
}

// AddField adds a field in to the schema.
// length は、その field が max 何 bytes 保存できるかの情報。VARCHAR(255) なら length は 255.
func (schema *Schema) AddField(fldname FieldName, typ FieldType, length int) {
	schema.fields = append(schema.fields, fldname)
	schema.info[fldname] = &FieldInfo{
		typ:    typ,
		length: length,
	}
}

// AddInt32Field adds an int field.
func (schema *Schema) AddInt32Field(fldname FieldName) {
	schema.AddField(fldname, Int32FieldType, 0)
}

// AddStringField adds an string field with maximum length is length.
// length is maximum length of string. It is not actual length of the value.
func (schema *Schema) AddStringField(fldname FieldName, length int) {
	schema.AddField(fldname, StringFieldType, length)
}

// Add adds other's field into the schema.
func (schema *Schema) Add(fldname FieldName, other *Schema) {
	typ := other.Type(fldname)

	length := other.Length(fldname)

	schema.AddField(fldname, typ, length)
}

// AddAllFields adds all fields of other into the schema.
func (schema *Schema) AddAllFields(other *Schema) {
	for _, fld := range other.fields {
		schema.Add(fld, other)
	}
}

// Fields returns schema fileds.
func (schema *Schema) Fields() []FieldName {
	return schema.fields
}

// Type returns field type.
func (schema *Schema) Type(fldname FieldName) FieldType {
	if v, found := schema.info[fldname]; found {
		return v.typ
	}

	return UnknownFieldType
}

// Length returns field byte length.
func (schema *Schema) Length(fldname FieldName) int {
	if v, found := schema.info[fldname]; found {
		return v.length
	}

	return -1
}

// Layout is model of table layout.
type Layout struct {
	schema   *Schema
	offsets  map[FieldName]int64
	slotsize int64
}

// NewLayout constructs Layout.
func NewLayout(schema *Schema) *Layout {
	pos := int64(RecordOffset) // flag for used/unused
	offsets := make(map[FieldName]int64)
	for _, fld := range schema.fields {
		offsets[fld] = pos

		// length in bytes
		switch schema.Type(fld) {
		case Int32FieldType:
			pos += common.Int32Length
		case StringFieldType:
			pos += common.Int32Length + int64(schema.Length(fld))
		case UnknownFieldType:
			log.Fatal(errors.New("Invalid field type"))
		}
	}

	return &Layout{
		schema:   schema,
		offsets:  offsets,
		slotsize: pos,
	}
}

// NewLayoutWithFields constructs a Layout with fields.
func NewLayoutWithFields(sch *Schema, offsets map[FieldName]int64, slotsize int64) *Layout {
	return &Layout{
		schema:   sch,
		offsets:  offsets,
		slotsize: slotsize,
	}
}

// Schema returns schema.
func (layout *Layout) Schema() *Schema {
	return layout.schema
}

// Offset returns field offset.
func (layout *Layout) Offset(fldname FieldName) int64 {
	return layout.offsets[fldname]
}

// SlotSize returns record slot size.
func (layout *Layout) SlotSize() int64 {
	return layout.slotsize
}

// Length returns byte size of given field name.
func (layout *Layout) Length(fldName FieldName) int {
	return layout.schema.Length(fldName)
}

// SlotCondition is flag of record slot condition.
type SlotCondition = int32

const (
	// Empty means the slot is empty.
	Empty SlotCondition = iota
	// Used means the slot is used.
	Used
)

// RecordPage is a model of RecordPage.
// Slot は record に usage flag をもたせたもの。
// Slot structure
// -------------------------------
// | usage flag (int32) | record |
// -------------------------------
//
// record structure
// ---------------------------
// | val 1 | val 2 | ...     |
// ---------------------------.
type RecordPage struct {
	txn    Transaction
	blk    Block
	layout *Layout
}

// NewRecordPage constructs a RecordPage.
func NewRecordPage(txn Transaction, blk Block, layout *Layout) (*RecordPage, error) {
	if err := txn.Pin(blk); err != nil {
		return nil, errors.Err(err, "Pin")
	}

	return &RecordPage{
		txn:    txn,
		blk:    blk,
		layout: layout,
	}, nil
}

// GetInt32 gets int32 from the block.
func (page *RecordPage) GetInt32(slotID SlotID, fldname FieldName) (int32, error) {
	offset := page.offset(slotID) + page.layout.Offset(fldname)

	return page.txn.GetInt32(page.blk, offset)
}

// SetInt32 sets int32 to the block.
func (page *RecordPage) SetInt32(slotID SlotID, fldname FieldName, val int32) error {
	offset := page.offset(slotID) + page.layout.Offset(fldname)

	return page.txn.SetInt32(page.blk, offset, val, true)
}

// GetString gets string from the block.
func (page *RecordPage) GetString(slotID SlotID, fldname FieldName) (string, error) {
	offset := page.offset(slotID) + page.layout.Offset(fldname)

	return page.txn.GetString(page.blk, offset)
}

// SetString sets the string from the block.
func (page *RecordPage) SetString(slotID SlotID, fldname FieldName, val string) error {
	offset := page.offset(slotID) + page.layout.Offset(fldname)

	return page.txn.SetString(page.blk, offset, val, true)
}

// Delete deletes the slot.
func (page *RecordPage) Delete(slotID SlotID) error {
	return page.setSlotCondition(slotID, Empty)
}

// Format formats blk.
func (page *RecordPage) Format() error {
	slotID := SlotID(0)
	for page.isValidSlot(slotID) {
		if err := page.txn.SetInt32(page.blk, page.offset(slotID), Empty, false); err != nil {
			return errors.Err(err, "SetInt32")
		}

		sch := page.layout.schema
		for _, fldname := range sch.fields {
			typ := sch.Type(fldname)
			fldpos := page.offset(slotID) + page.layout.Offset(fldname)
			switch typ {
			case Int32FieldType:
				if err := page.txn.SetInt32(page.blk, fldpos, 0, false); err != nil {
					return errors.Err(err, "SetInt32")
				}
			case StringFieldType:
				if err := page.txn.SetString(page.blk, fldpos, "", false); err != nil {
					return errors.Err(err, "SetString")
				}
			case UnknownFieldType:
				log.Fatal(errors.New("unexpected record type"))
			}
		}
		slotID++
	}

	return nil
}

// NextUsedSlot returns the slot id with Used flag after slot.
// 引数に与えた SlotID よりもあとにある used slot の ID を返却する。
func (page *RecordPage) NextUsedSlot(slotID SlotID) (SlotID, error) {
	return page.searchAfter(slotID, Used)
}

// InsertAfter searches the slot id after slot with Empty flag, set Used flag and returns its id.
// 与えられた slotID よりもあとにある empty flag の slot を used にし、その ID を返却する。
// MEMO: empty を探す作業と、flag を used にする作業は method を分けたほうがよいのではないか？
func (page *RecordPage) InsertAfter(slotID SlotID) (SlotID, error) {
	newSlot, err := page.searchAfter(slotID, Empty)
	if err != nil {
		return 0, errors.Err(err, "searchAfter")
	}
	if newSlot >= 0 {
		err := page.setSlotCondition(newSlot, Used)
		if err != nil {
			return 0, errors.Err(err, "setSlotCondition")
		}
	}

	return newSlot, nil
}

// searchAfter searches slot id with given flag after slot.
func (page *RecordPage) searchAfter(slotID SlotID, flag SlotCondition) (SlotID, error) {
	slotID++
	for page.isValidSlot(slotID) {
		currFlag, err := page.GetSlotCondition(slotID)
		if err != nil {
			return 0, errors.Wrap(err, "GetSlotCondition")
		}
		if flag == currFlag {
			return slotID, nil
		}
		slotID++
	}

	return -1, nil
}

func (page *RecordPage) isValidSlot(slotID SlotID) bool {
	off := page.offset(slotID + 1)
	x := int64(page.txn.BlockSize())

	return off <= x
}

func (page *RecordPage) setSlotCondition(slotID SlotID, flag SlotCondition) error {
	return page.txn.SetInt32(page.blk, page.offset(slotID), flag, true)
}

// GetSlotCondition gets slot condition (used/unused).
func (page *RecordPage) GetSlotCondition(slotID SlotID) (SlotCondition, error) {
	return page.txn.GetInt32(page.blk, page.offset(slotID))
}

func (page *RecordPage) offset(slotID SlotID) int64 {
	return int64(slotID) * page.layout.slotsize
}

// Block returns RecordPage's block.
func (page *RecordPage) Block() Block {
	return page.blk
}

// RecordID is identifier of record.
type RecordID struct {
	blkNum BlockNumber
	slotID SlotID
}

// NewRecordID constructs RecordID.
func NewRecordID(blkNum BlockNumber, slotID SlotID) RecordID {
	return RecordID{
		blkNum: blkNum,
		slotID: slotID,
	}
}

// BlockNumber returns block number.
func (rid RecordID) BlockNumber() BlockNumber {
	return rid.blkNum
}

// SlotID returns slot id.
func (rid RecordID) SlotID() SlotID {
	return rid.slotID
}

// Equal checks equality of rid and other.
func (rid RecordID) Equal(other RecordID) bool {
	return rid.blkNum == other.blkNum && rid.slotID == other.slotID
}
