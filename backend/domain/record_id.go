package domain

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

// NewZeroRecordID returns zero value of record id.
func NewZeroRecordID() RecordID {
	return RecordID{}
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
