package logrecord

import (
	"github.com/goropikari/simpledbgo/backend/tx/logrecord/protobuf"
	"github.com/goropikari/simpledbgo/domain"
	"google.golang.org/protobuf/proto"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

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

	// Checkpoint is checkpoint record type.
	Checkpoint

	// SetInt32 is set int32 record type.
	SetInt32

	// SetString is set string record type.
	SetString
)

// TxVisitor is an interface of visitor.
type TxVisitor interface {
	Pin(domain.Block) error
	Unpin(domain.Block)
	UndoSetInt32(*SetInt32Record) error
	UndoSetString(*SetStringRecord) error
}

// LogRecorder is an interface of log record.
type LogRecorder interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Operator() RecordType
	TxNumber() domain.TransactionNumber
	Undo(TxVisitor) error
}

type baseRecord struct{}

// Undo is dummy method for implementing LogRecorder.
func (rec *baseRecord) Undo(visitor TxVisitor) error {
	return nil
}

// StartRecord is a model of start log record.
type StartRecord struct {
	baseRecord
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

// TxNumber returns the transaction number.
func (rec *StartRecord) TxNumber() domain.TransactionNumber {
	return rec.TxNum
}

// CommitRecord is a model of commit log record.
type CommitRecord struct {
	baseRecord
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

// TxNumber returns the transaction number.
func (rec *CommitRecord) TxNumber() domain.TransactionNumber {
	return rec.TxNum
}

// RollbackRecord is a model of commit log record.
type RollbackRecord struct {
	baseRecord
	TxNum domain.TransactionNumber
}

// Unmarshal parses the proto message in b and places the result in rec.
func (rec *RollbackRecord) Unmarshal(b []byte) error {
	pb := &protobuf.RollbackRecord{}
	if err := proto.Unmarshal(b, pb); err != nil {
		return err
	}

	rec.TxNum = domain.TransactionNumber(pb.Txnum)

	return nil
}

// Marshal encodes the rec.
func (rec *RollbackRecord) Marshal() ([]byte, error) {
	pb := &protobuf.RollbackRecord{
		Txnum: int32(rec.TxNum),
	}

	return proto.Marshal(pb)
}

// Operator returns Commit.
func (rec *RollbackRecord) Operator() RecordType {
	return Rollback
}

// TxNumber returns the transaction number.
func (rec *RollbackRecord) TxNumber() domain.TransactionNumber {
	return rec.TxNum
}

// CheckpointRecord is a model of commit log record.
type CheckpointRecord struct {
	baseRecord
}

// Unmarshal parses the proto message in b and places the result in rec.
func (rec *CheckpointRecord) Unmarshal(b []byte) error {
	pb := &protobuf.CheckpointRecord{}
	if err := proto.Unmarshal(b, pb); err != nil {
		return err
	}

	return nil
}

// Marshal encodes the rec.
func (rec *CheckpointRecord) Marshal() ([]byte, error) {
	pb := &protobuf.CheckpointRecord{}

	return proto.Marshal(pb)
}

// Operator returns Commit.
func (rec *CheckpointRecord) Operator() RecordType {
	return Checkpoint
}

// TxNumber returns the transaction number.
func (rec *CheckpointRecord) TxNumber() domain.TransactionNumber {
	return domain.DummyTransactionNumber
}

// SetInt32Record is a model of set int32 log record.
type SetInt32Record struct {
	baseRecord
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

// TxNumber returns the transaction number.
func (rec *SetInt32Record) TxNumber() domain.TransactionNumber {
	return rec.TxNum
}

// Undo undoes set int32 operation.
func (rec *SetInt32Record) Undo(visitor TxVisitor) error {
	return visitor.UndoSetInt32(rec)
}

// SetStringRecord is a model of set string log record.
type SetStringRecord struct {
	baseRecord
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

// TxNumber returns the transaction number.
func (rec *SetStringRecord) TxNumber() domain.TransactionNumber {
	return rec.TxNum
}

// Undo undoes set string operation.
func (rec *SetStringRecord) Undo(visitor TxVisitor) error {
	return visitor.UndoSetString(rec)
}
