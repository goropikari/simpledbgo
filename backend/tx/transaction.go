package tx

import (
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/tx/logrecord"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/meta"
)

// RecordType is type of log record.
type RecordType = int32

const (
	// Start is start record type.
	Start RecordType = iota

	// Commit is commit record type.
	Commit

	// SetInt32 is set int32 record type.
	SetInt32

	// SetString is set string record type.
	SetString
)

// Transaction is a model of transaction.
type Transaction struct {
	fileMgr    domain.FileManager
	logMgr     domain.LogManager
	bufferMgr  domain.BufferManager
	concurMgr  *ConcurrencyManager
	bufferList *BufferList
	number     domain.TransactionNumber
}

// NewTransaction constructs Transaction.
func NewTransaction(fileMgr domain.FileManager, logMgr domain.LogManager, bufferMgr domain.BufferManager, concurMgr *ConcurrencyManager, gen *NumberGenerator) (*Transaction, error) {
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

		lsn, err = tx.writeSetInt32Log(*buf.Block(), offset, oldval)
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
		// recoveryMgr.setString
		oldval, err := buf.Page().GetString(offset)
		if err != nil {
			return err
		}
		lsn, err = tx.writeSetStringLog(*buf.Block(), offset, oldval)
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

	return tx.writeLog(Start, record)
}

func (tx *Transaction) writeCommitLog() (domain.LSN, error) {
	record := &logrecord.CommitRecord{TxNum: tx.number}

	return tx.writeLog(Commit, record)
}

func (tx *Transaction) writeSetInt32Log(blk domain.Block, offset int64, val int32) (domain.LSN, error) {
	record := &logrecord.SetInt32Record{
		FileName:    blk.FileName(),
		TxNum:       tx.number,
		BlockNumber: blk.Number(),
		Offset:      offset,
		Val:         val,
	}

	return tx.writeLog(SetInt32, record)
}

func (tx *Transaction) writeSetStringLog(blk domain.Block, offset int64, val string) (domain.LSN, error) {
	record := &logrecord.SetStringRecord{
		FileName:    blk.FileName(),
		TxNum:       tx.number,
		BlockNumber: blk.Number(),
		Offset:      offset,
		Val:         val,
	}

	return tx.writeLog(SetString, record)
}

func (tx *Transaction) writeLog(typ RecordType, record logrecord.LogRecorder) (domain.LSN, error) {
	data, err := record.Marshal()
	if err != nil {
		return domain.DummyLSN, err
	}

	buf := bytes.NewBuffer(meta.Int32Length*2 + len(data))
	if err := buf.SetInt32(0, typ); err != nil {
		return domain.DummyLSN, err
	}

	if err := buf.SetBytes(meta.Int32Length, data); err != nil {
		return domain.DummyLSN, err
	}

	lsn, err := tx.logMgr.AppendRecord(buf.GetData())
	if err != nil {
		return domain.DummyLSN, err
	}

	return lsn, nil
}
