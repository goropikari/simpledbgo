package logrecord

import (
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx/logrecord/protobuf"
	"google.golang.org/protobuf/proto"
)

// RecordType is type of log record.
type RecordType = int32

const (
	// Unknown is an unknown record type.
	Unknown RecordType = iota

	// Start is start record type.
	Start

	// Commit is commit record type.
	Commit

	// Rollback is rollback record type.
	Rollback

	// SetInt32 is set int32 record type.
	SetInt32

	// SetString is set string record type.
	SetString
)

// LogRecorder is an interface of log record.
type LogRecorder interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Operator() RecordType
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

// Operator returns Start.
func (rec *StartRecord) Operator() RecordType {
	return Start
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

// Operator returns Commit.
func (rec *CommitRecord) Operator() RecordType {
	return Commit
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

// Operator returns SetInt32.
func (rec *SetInt32Record) Operator() RecordType {
	return SetInt32
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

// Operator returns SetString.
func (rec *SetStringRecord) Operator() RecordType {
	return SetString
}
