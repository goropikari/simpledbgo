package logrecord

import (
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx/logrecord/protobuf"
	"google.golang.org/protobuf/proto"
)

// SetInt32Record is a model of set int32 log record.
type SetInt32Record struct {
	FileName    domain.FileName
	TxNum       domain.TransactionNumber
	BlockNumber domain.BlockNumber
	Offset      int64
	Val         int32
}

// Unmarshal parses the proto message in b and places the result in rec.
func (rec *SetInt32Record) Unmarshal(b []byte) error {
	pb := &protobuf.SetInt32Record{}
	if err := proto.Unmarshal(b, pb); err != nil {
		return err
	}

	var err error
	rec.FileName, err = domain.NewFileName(pb.Filename)
	if err != nil {
		return err
	}

	rec.TxNum = domain.TransactionNumber(pb.Txnum)
	rec.BlockNumber, err = domain.NewBlockNumber(pb.BlockNumber)
	if err != nil {
		return err
	}

	rec.Offset = pb.Offset
	rec.Val = pb.Val

	return nil
}

// Marshal encodes the rec.
func (rec *SetInt32Record) Marshal() ([]byte, error) {
	pb := &protobuf.SetInt32Record{
		Filename:    rec.FileName.String(),
		Txnum:       int32(rec.TxNum),
		BlockNumber: int32(rec.BlockNumber),
		Offset:      rec.Offset,
		Val:         rec.Val,
	}

	return proto.Marshal(pb)
}
