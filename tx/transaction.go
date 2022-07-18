package tx

import (
	"github.com/goropikari/simpledbgo/common"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/tx/logrecord"
)

// Transaction is a model of transaction.
type Transaction struct {
	fileMgr    domain.FileManager
	logMgr     domain.LogManager
	bufferMgr  domain.BufferPoolManager
	concurMgr  *ConcurrencyManager
	bufferList *BufferList
	number     domain.TransactionNumber
}

// NewTransaction constructs Transaction.
func NewTransaction(fileMgr domain.FileManager, logMgr domain.LogManager, bufferMgr domain.BufferPoolManager, lt *LockTable, gen domain.TxNumberGenerator) (*Transaction, error) {
	txn := &Transaction{
		fileMgr:    fileMgr,
		logMgr:     logMgr,
		bufferMgr:  bufferMgr,
		concurMgr:  NewConcurrencyManager(lt),
		bufferList: NewBufferList(bufferMgr),
		number:     gen.Generate(),
	}

	if _, err := txn.writeStartLog(); err != nil {
		return nil, errors.Err(err, "writeStartLog")
	}

	return txn, nil
}

// Pin pins the blk by tx.
func (tx *Transaction) Pin(blk domain.Block) error {
	if err := tx.bufferList.Pin(blk); err != nil {
		return errors.Err(err, "Pin")
	}

	return nil
}

// Unpin unpins the blk by tx.
func (tx *Transaction) Unpin(blk domain.Block) {
	tx.bufferList.Unpin(blk)
}

// Commit commits the transaction.
func (tx *Transaction) Commit() error {
	if err := tx.commit(); err != nil {
		return errors.Err(err, "commit")
	}

	tx.concurMgr.Release()
	tx.bufferList.UnpinAll()

	return nil
}

func (tx *Transaction) commit() error {
	if err := tx.bufferMgr.FlushAll(tx.number); err != nil {
		return errors.Err(err, "FlushAll")
	}

	lsn, err := tx.writeCommitLog()
	if err != nil {
		return errors.Err(err, "writeCommitLog")
	}

	if err := tx.logMgr.FlushLSN(lsn); err != nil {
		return errors.Err(err, "FlushLSN")
	}

	return nil
}

// Rollback rollbacks the transaction.
func (tx *Transaction) Rollback() error {
	if err := tx.rollback(); err != nil {
		return errors.Err(err, "rollback")
	}

	tx.concurMgr.Release()
	tx.bufferList.UnpinAll()

	return nil
}

func (tx *Transaction) rollback() error {
	iter, err := tx.logMgr.Iterator()
	if err != nil {
		return errors.Err(err, "Iterator")
	}

	for iter.HasNext() {
		data, err := iter.Next()
		if err != nil {
			return errors.Err(err, "Next")
		}

		record, err := ParseRecord(data)
		if err != nil {
			return errors.Err(err, "ParseRecord")
		}

		if record.TxNumber() == tx.number {
			if record.Operator() == logrecord.Start {
				break
			}
			if err := record.Undo(tx); err != nil {
				return errors.Err(err, "Undo")
			}
		}
	}
	if err := iter.Err(); err != nil {
		return errors.Err(err, "HasNext")
	}

	if err := tx.bufferMgr.FlushAll(tx.number); err != nil {
		return errors.Err(err, "FlushAll")
	}

	lsn, err := tx.writeRollbackLog()
	if err != nil {
		return errors.Err(err, "writeRollbackLog")
	}

	if err := tx.logMgr.FlushLSN(lsn); err != nil {
		return errors.Err(err, "FlushLSN")
	}

	return nil
}

// Recover recovers a database.
func (tx *Transaction) Recover() error {
	if err := tx.bufferMgr.FlushAll(tx.number); err != nil {
		return errors.Err(err, "FlushAll")
	}

	if err := tx.recover(); err != nil {
		return errors.Err(err, "recover")
	}

	return nil
}

func (tx *Transaction) recover() error {
	finishedTxns := make(map[domain.TransactionNumber]bool)
	iter, err := tx.logMgr.Iterator()
	if err != nil {
		return errors.Err(err, "Iterator")
	}

	for iter.HasNext() {
		data, err := iter.Next()
		if err != nil {
			return errors.Err(err, "Next")
		}

		record, err := ParseRecord(data)
		if err != nil {
			return errors.Err(err, "ParseRecord")
		}

		op := record.Operator()
		if op == logrecord.Checkpoint {
			break
		}

		if op == logrecord.Commit || op == logrecord.Rollback {
			finishedTxns[record.TxNumber()] = true
		} else if _, found := finishedTxns[record.TxNumber()]; !found {
			if err := record.Undo(tx); err != nil {
				return errors.Err(err, "Undo")
			}
		}
	}
	if err := iter.Err(); err != nil {
		return errors.Err(err, "HasNext")
	}

	if err := tx.bufferMgr.FlushAll(tx.number); err != nil {
		return errors.Err(err, "FlushAll")
	}

	lsn, err := tx.writeCheckpointLog()
	if err != nil {
		return errors.Err(err, "writeCheckpointLog")
	}

	if err := tx.logMgr.FlushLSN(lsn); err != nil {
		return errors.Err(err, "FlushLSN")
	}

	return nil
}

// UndoSetInt32 undoes SetInt32 operation.
func (tx *Transaction) UndoSetInt32(rec *logrecord.SetInt32Record) error {
	blk := domain.NewBlock(rec.FileName, rec.BlockNumber)

	if err := tx.Pin(blk); err != nil {
		return errors.Err(err, "Pin")
	}

	if err := tx.SetInt32(blk, rec.Offset, rec.Val, false); err != nil {
		return errors.Err(err, "SetInt32")
	}

	tx.Unpin(blk)

	return nil
}

// UndoSetString undoes SetString.
func (tx *Transaction) UndoSetString(rec *logrecord.SetStringRecord) error {
	blk := domain.NewBlock(rec.FileName, rec.BlockNumber)

	if err := tx.Pin(blk); err != nil {
		return errors.Err(err, "Pin")
	}

	if err := tx.SetString(blk, rec.Offset, rec.Val, false); err != nil {
		return errors.Err(err, "SetString")
	}

	tx.Unpin(blk)

	return nil
}

// GetInt32 gets int32 from the blk at offset.
func (tx *Transaction) GetInt32(blk domain.Block, offset int64) (int32, error) {
	if err := tx.concurMgr.SLock(blk); err != nil {
		return 0, errors.Err(err, "SLock")
	}

	buf := tx.bufferList.GetBuffer(blk)
	x, err := buf.Page().GetInt32(offset)
	if err != nil {
		return 0, errors.Err(err, "GetInt32")
	}

	return x, nil
}

// SetInt32 sets int32 on the given block.
func (tx *Transaction) SetInt32(blk domain.Block, offset int64, val int32, writeLog bool) error {
	if err := tx.concurMgr.XLock(blk); err != nil {
		return errors.Err(err, "XLock")
	}

	buf := tx.bufferList.GetBuffer(blk)
	lsn := domain.DummyLSN
	if writeLog {
		var err error
		oldval, err := buf.Page().GetInt32(offset)
		if err != nil {
			return errors.Err(err, "GetInt32")
		}

		lsn, err = tx.writeSetInt32Log(buf.Block(), offset, oldval)
		if err != nil {
			return errors.Err(err, "writeSetInt32Log")
		}
	}

	if err := buf.Page().SetInt32(offset, val); err != nil {
		return errors.Err(err, "SetInt32")
	}

	buf.SetModifiedTxNumber(tx.number, lsn)

	return nil
}

// GetString gets string from the blk.
func (tx *Transaction) GetString(blk domain.Block, offset int64) (string, error) {
	if err := tx.concurMgr.SLock(blk); err != nil {
		return "", errors.Err(err, "SLock")
	}

	buf := tx.bufferList.GetBuffer(blk)

	return buf.Page().GetString(offset)
}

// SetString sets string on the blk.
func (tx *Transaction) SetString(blk domain.Block, offset int64, val string, writeLog bool) error {
	if err := tx.concurMgr.XLock(blk); err != nil {
		return errors.Err(err, "XLock")
	}

	buf := tx.bufferList.GetBuffer(blk)
	lsn := domain.DummyLSN
	if writeLog {
		oldval, err := buf.Page().GetString(offset)
		if err != nil {
			return errors.Err(err, "GetString")
		}
		lsn, err = tx.writeSetStringLog(buf.Block(), offset, oldval)
		if err != nil {
			return errors.Err(err, "writeSetStringLog")
		}
	}

	page := buf.Page()
	if err := page.SetString(offset, val); err != nil {
		return errors.Err(err, "SetString")
	}

	buf.SetModifiedTxNumber(tx.number, lsn)

	return nil
}

func (tx *Transaction) writeStartLog() (domain.LSN, error) {
	record := &logrecord.StartRecord{
		TxNum: tx.number,
	}

	return tx.writeLog(logrecord.Start, record)
}

func (tx *Transaction) writeCommitLog() (domain.LSN, error) {
	record := &logrecord.CommitRecord{TxNum: tx.number}

	return tx.writeLog(logrecord.Commit, record)
}

func (tx *Transaction) writeCheckpointLog() (domain.LSN, error) {
	record := &logrecord.CheckpointRecord{}

	return tx.writeLog(logrecord.Checkpoint, record)
}

func (tx *Transaction) writeSetInt32Log(blk domain.Block, offset int64, val int32) (domain.LSN, error) {
	record := &logrecord.SetInt32Record{
		FileName:    blk.FileName(),
		TxNum:       tx.number,
		BlockNumber: blk.Number(),
		Offset:      offset,
		Val:         val,
	}

	return tx.writeLog(logrecord.SetInt32, record)
}

func (tx *Transaction) writeSetStringLog(blk domain.Block, offset int64, val string) (domain.LSN, error) {
	record := &logrecord.SetStringRecord{
		FileName:    blk.FileName(),
		TxNum:       tx.number,
		BlockNumber: blk.Number(),
		Offset:      offset,
		Val:         val,
	}

	return tx.writeLog(logrecord.SetString, record)
}

func (tx *Transaction) writeRollbackLog() (domain.LSN, error) {
	record := &logrecord.RollbackRecord{TxNum: tx.number}

	return tx.writeLog(logrecord.Rollback, record)
}

// log record structure
// ----------------------------------------------------------------------
// | log record type (int32) | record length (uint32) | record (varlen) |
// ----------------------------------------------------------------------.
func (tx *Transaction) writeLog(typ logrecord.RecordType, record logrecord.LogRecorder) (domain.LSN, error) {
	data, err := record.Marshal()
	if err != nil {
		return domain.DummyLSN, errors.Err(err, "Marshal")
	}

	buf := bytes.NewBuffer(common.Int32Length + common.Uint32Length + len(data))
	if err := buf.SetInt32(0, typ); err != nil {
		return domain.DummyLSN, errors.Err(err, "SetInt32")
	}

	if err := buf.SetBytes(common.Int32Length, data); err != nil {
		return domain.DummyLSN, errors.Err(err, "SetBytes")
	}

	lsn, err := tx.logMgr.AppendRecord(buf.GetData())
	if err != nil {
		return domain.DummyLSN, errors.Err(err, "AppendRecord")
	}

	return lsn, nil
}

// BlockLength returns block length of the `filename`.
// original method name is `size`.
func (tx *Transaction) BlockLength(filename domain.FileName) (int32, error) {
	dummyBlk := domain.NewDummyBlock(filename)
	if err := tx.concurMgr.SLock(dummyBlk); err != nil {
		return 0, errors.Err(err, "SLock")
	}

	return tx.fileMgr.BlockLength(filename)
}

// ExtendFile extends the file by a block.
func (tx *Transaction) ExtendFile(filename domain.FileName) (domain.Block, error) {
	dummyBlk := domain.NewDummyBlock(filename)
	if err := tx.concurMgr.XLock(dummyBlk); err != nil {
		return domain.Block{}, errors.Err(err, "XLock")
	}

	return tx.fileMgr.ExtendFile(filename)
}

// BlockSize returns block size.
func (tx *Transaction) BlockSize() domain.BlockSize {
	return tx.fileMgr.BlockSize()
}

// Available returns the number of available buffers.
func (tx *Transaction) Available() int {
	return tx.bufferMgr.Available()
}

func (tx *Transaction) GetTxNum() domain.TransactionNumber {
	return tx.number
}
