package domain

// RecordID is identifier of record.
type RecordID struct {
	blkNum BlockNumber
	slotID SlotID
}

// NewRecordID constructs RecordID.
func NewRecordID(blkNum BlockNumber, slot SlotID) RecordID {
	return RecordID{
		blkNum: blkNum,
		slotID: slot,
	}
}

// BlockNumber returns block number.
func (rid RecordID) BlockNumber() BlockNumber {
	return rid.blkNum
}

// // SlotID returns slot id.
// func (rid RecordID) SlotID() SlotID {
// 	return rid.slotID
// }
