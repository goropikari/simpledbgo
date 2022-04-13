package logrecord

import (
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx/logrecord/protobuf"
	"google.golang.org/protobuf/proto"
)

// LogRecorder is an interface of log record.
type LogRecorder interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
}

// StartRecord is a model of start log record.
type StartRecord struct {
	TxNum domain.TransactionNumber
}

// Unmarshal parses the proto message in b and places the result in rec.
func (rec *StartRecord) Unmarshal(b []byte) error {
	pb := &protobuf.StartRecord{}
	if err := proto.Unmarshal(b, pb); err != nil {
		return err
	}

	rec.TxNum = domain.TransactionNumber(pb.Txnum)

	return nil
}

// Marshal encodes the rec.
func (rec *StartRecord) Marshal() ([]byte, error) {
	pb := &protobuf.StartRecord{
		Txnum: int32(rec.TxNum),
	}

	return proto.Marshal(pb)
}

// CommitRecord is a model of commit log record.
type CommitRecord struct {
	TxNum domain.TransactionNumber
}

// Unmarshal parses the proto message in b and places the result in rec.
func (rec *CommitRecord) Unmarshal(b []byte) error {
	pb := &protobuf.CommitRecord{}
	if err := proto.Unmarshal(b, pb); err != nil {
		return err
	}

	rec.TxNum = domain.TransactionNumber(pb.Txnum)

	return nil
}

// Marshal encodes the rec.
func (rec *CommitRecord) Marshal() ([]byte, error) {
	pb := &protobuf.CommitRecord{
		Txnum: int32(rec.TxNum),
	}

	return proto.Marshal(pb)
}

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

// SetStringRecord is a model of set string log record.
type SetStringRecord struct {
	FileName    domain.FileName
	TxNum       domain.TransactionNumber
	BlockNumber domain.BlockNumber
	Offset      int64
	Val         string
}

// Unmarshal parses the proto message in b and places the result in rec.
func (rec *SetStringRecord) Unmarshal(b []byte) error {
	pb := &protobuf.SetStringRecord{}
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
func (rec *SetStringRecord) Marshal() ([]byte, error) {
	pb := &protobuf.SetStringRecord{
		Filename:    rec.FileName.String(),
		Txnum:       int32(rec.TxNum),
		BlockNumber: int32(rec.BlockNumber),
		Offset:      rec.Offset,
		Val:         rec.Val,
	}

	return proto.Marshal(pb)
}
