package tx

import (
	"github.com/goropikari/simpledbgo/common"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/tx/logrecord"
)

// Transaction is a model of transaction.
type Transaction struct {
	fileMgr    domain.FileManager
	logMgr     domain.LogManager
	bufferMgr  domain.BufferManager
	concurMgr  domain.ConcurrencyManager
	bufferList *BufferList
	number     domain.TransactionNumber
}

// NewTransaction constructs Transaction.
func NewTransaction(fileMgr domain.FileManager, logMgr domain.LogManager, bufferMgr domain.BufferManager, concurMgr domain.ConcurrencyManager, gen domain.TxNumberGenerator) (*Transaction, error) {
	txn := &Transaction{
		fileMgr:    fileMgr,
		logMgr:     logMgr,
		bufferMgr:  bufferMgr,
		concurMgr:  concurMgr,
		bufferList: NewBufferList(bufferMgr),
		number:     gen.Generate(),
	}

	if _, err := txn.writeStartLog(); err != nil {
		return nil, err
	}

	return txn, nil
}

// Pin pins the blk by tx.
func (tx *Transaction) Pin(blk domain.Block) error {
	if err := tx.bufferList.Pin(blk); err != nil {
		return err
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
		return err
	}

	tx.concurMgr.Release()
	tx.bufferList.UnpinAll()

	return nil
}

func (tx *Transaction) commit() error {
	if err := tx.bufferMgr.FlushAll(tx.number); err != nil {
		return err
	}

	lsn, err := tx.writeCommitLog()
	if err != nil {
		return err
	}

	if err := tx.logMgr.FlushLSN(lsn); err != nil {
		return err
	}

	return nil
}

// Rollback rollbacks the transaction.
func (tx *Transaction) Rollback() error {
	if err := tx.rollback(); err != nil {
		return err
	}

	tx.concurMgr.Release()
	tx.bufferList.UnpinAll()

	return nil
}

func (tx *Transaction) rollback() error {
	iter, err := tx.logMgr.Iterator()
	if err != nil {
		return err
	}

	for iter.HasNext() {
		data, err := iter.Next()
		if err != nil {
			return err
		}

		record, err := ParseRecord(data)
		if err != nil {
			return err
		}

		if record.TxNumber() == tx.number {
			if record.Operator() == logrecord.Start {
				break
			}
			if err := record.Undo(tx); err != nil {
				return err
			}
		}
	}

	if err := tx.bufferMgr.FlushAll(tx.number); err != nil {
		return err
	}

	lsn, err := tx.writeRollbackLog()
	if err != nil {
		return err
	}

	if err := tx.logMgr.FlushLSN(lsn); err != nil {
		return err
	}

	return nil
}

// Recover recovers a database.
func (tx *Transaction) Recover() error {
	if err := tx.bufferMgr.FlushAll(tx.number); err != nil {
		return err
	}

	if err := tx.recover(); err != nil {
		return err
	}

	return nil
}

func (tx *Transaction) recover() error {
	finishedTxns := make(map[domain.TransactionNumber]bool)
	iter, err := tx.logMgr.Iterator()
	if err != nil {
		return err
	}

	for iter.HasNext() {
		data, err := iter.Next()
		if err != nil {
			return err
		}

		record, err := ParseRecord(data)
		if err != nil {
			return err
		}

		op := record.Operator()
		if op == logrecord.Checkpoint {
			break
		}

		if op == logrecord.Commit || op == logrecord.Rollback {
			finishedTxns[record.TxNumber()] = true
		} else if _, found := finishedTxns[record.TxNumber()]; !found {
			if err := record.Undo(tx); err != nil {
				return err
			}
		}
	}

	if err := tx.bufferMgr.FlushAll(tx.number); err != nil {
		return err
	}

	lsn, err := tx.writeCheckpointLog()
	if err != nil {
		return err
	}

	if err := tx.logMgr.FlushLSN(lsn); err != nil {
		return err
	}

	return nil
}

// UndoSetInt32 undoes SetInt32 operation.
func (tx *Transaction) UndoSetInt32(rec *logrecord.SetInt32Record) error {
	blk := domain.NewBlock(rec.FileName, rec.BlockNumber)

	if err := tx.Pin(blk); err != nil {
		return err
	}

	if err := tx.SetInt32(blk, rec.Offset, rec.Val, false); err != nil {
		return err
	}

	tx.Unpin(blk)

	return nil
}

// UndoSetString undoes SetString.
func (tx *Transaction) UndoSetString(rec *logrecord.SetStringRecord) error {
	blk := domain.NewBlock(rec.FileName, rec.BlockNumber)

	if err := tx.Pin(blk); err != nil {
		return err
	}

	if err := tx.SetString(blk, rec.Offset, rec.Val, false); err != nil {
		return err
	}

	tx.Unpin(blk)

	return nil
}

// GetInt32 gets int32 from the blk at offset.
func (tx *Transaction) GetInt32(blk domain.Block, offset int64) (int32, error) {
	if err := tx.concurMgr.SLock(blk); err != nil {
		return 0, err
	}

	buf := tx.bufferList.GetBuffer(blk)
	x, err := buf.Page().GetInt32(offset)
	if err != nil {
		return 0, err
	}

	return x, nil
}

// SetInt32 sets int32 on the given block.
func (tx *Transaction) SetInt32(blk domain.Block, offset int64, val int32, writeLog bool) error {
	if err := tx.concurMgr.XLock(blk); err != nil {
		return err
	}

	buf := tx.bufferList.GetBuffer(blk)
	lsn := domain.DummyLSN
	if writeLog {
		var err error
		oldval, err := buf.Page().GetInt32(offset)
		if err != nil {
			return err
		}

		lsn, err = tx.writeSetInt32Log(buf.Block(), offset, oldval)
		if err != nil {
			return err
		}
	}

	if err := buf.Page().SetInt32(offset, val); err != nil {
		return err
	}

	buf.SetModifiedTxNumber(tx.number, lsn)

	return nil
}

// GetString gets string from the blk.
func (tx *Transaction) GetString(blk domain.Block, offset int64) (string, error) {
	if err := tx.concurMgr.SLock(blk); err != nil {
		return "", err
	}

	buf := tx.bufferList.GetBuffer(blk)

	return buf.Page().GetString(offset)
}

// SetString sets string on the blk.
func (tx *Transaction) SetString(blk domain.Block, offset int64, val string, writeLog bool) error {
	if err := tx.concurMgr.XLock(blk); err != nil {
		return err
	}

	buf := tx.bufferList.GetBuffer(blk)
	lsn := domain.DummyLSN
	if writeLog {
		oldval, err := buf.Page().GetString(offset)
		if err != nil {
			return err
		}
		lsn, err = tx.writeSetStringLog(buf.Block(), offset, oldval)
		if err != nil {
			return err
		}
	}

	page := buf.Page()
	if err := page.SetString(offset, val); err != nil {
		return err
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

func (tx *Transaction) writeLog(typ logrecord.RecordType, record logrecord.LogRecorder) (domain.LSN, error) {
	data, err := record.Marshal()
	if err != nil {
		return domain.DummyLSN, err
	}

	buf := bytes.NewBuffer(common.Int32Length*2 + len(data))
	if err := buf.SetInt32(0, typ); err != nil {
		return domain.DummyLSN, err
	}

	if err := buf.SetBytes(common.Int32Length, data); err != nil {
		return domain.DummyLSN, err
	}

	lsn, err := tx.logMgr.AppendRecord(buf.GetData())
	if err != nil {
		return domain.DummyLSN, err
	}

	return lsn, nil
}

// BlockLength returns block length of the `filename`.
func (tx *Transaction) BlockLength(filename domain.FileName) (int32, error) {
	dummyBlk := domain.NewDummyBlock(filename)
	if err := tx.concurMgr.SLock(dummyBlk); err != nil {
		return 0, err
	}

	return tx.fileMgr.BlockLength(filename)
}

// ExtendFile extends the file by a block.
func (tx *Transaction) ExtendFile(filename domain.FileName) (domain.Block, error) {
	dummyBlk := domain.NewDummyBlock(filename)
	if err := tx.concurMgr.XLock(dummyBlk); err != nil {
		return domain.Block{}, err
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
