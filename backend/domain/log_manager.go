package domain

import (
	"errors"
	"sync"
)

const int32Length = 4

// LogManagerConfig is a configuration of log manager.
type LogManagerConfig struct {
	LogFileName string
}

// LogManager is a log manager.
type LogManager struct {
	mu           sync.Mutex
	fileMgr      *FileManager
	logFileName  FileName
	currentBlock *Block
	logPage      *Page
	latestLSN    int32
	lastSavedLSN int32
}

// NewLogManager is a constructor of LogManager.
func NewLogManager(fileMgr *FileManager, block *Block, page *Page, config LogManagerConfig) (*LogManager, error) {
	logFileName, err := NewFileName(config.LogFileName)
	if err != nil {
		return nil, err
	}

	return &LogManager{
		mu:           sync.Mutex{},
		fileMgr:      fileMgr,
		logFileName:  logFileName,
		currentBlock: block,
		logPage:      page,
		latestLSN:    0,
		lastSavedLSN: 0,
	}, nil
}

// PrepareLogManager prepares a block and a page for initializing LogManager.
// If given file is empty, extend a file by block size.
func PrepareLogManager(fileMgr *FileManager, factory *PageFactory, fileName FileName) (*Block, *Page, error) {
	page, err := factory.Create()
	if err != nil {
		return nil, nil, err
	}

	blklen, err := fileMgr.BlockLength(fileName)
	if err != nil {
		return nil, nil, err
	}

	if blklen == 0 {
		blk, err := fileMgr.ExtendFile(fileName)
		if err != nil {
			return nil, nil, err
		}

		err = page.SetInt32(0, int32(fileMgr.BlockSize()))
		if err != nil {
			return nil, nil, err
		}

		err = fileMgr.CopyPageToBlock(page, blk)
		if err != nil {
			return nil, nil, err
		}

		return blk, page, nil
	}

	blknum, err := NewBlockNumber(blklen - 1)
	blk := NewBlock(fileName, fileMgr.BlockSize(), blknum)

	err = fileMgr.CopyBlockToPage(blk, page)
	if err != nil {
		return nil, nil, err
	}

	return blk, page, nil
}

// FlushLSN flushes by lsn.
func (mgr *LogManager) FlushLSN(lsn int32) error {
	if lsn >= mgr.lastSavedLSN {
		return mgr.Flush()
	}

	return nil
}

// Flush flushes the log page.
func (mgr *LogManager) Flush() error {
	err := mgr.fileMgr.CopyPageToBlock(mgr.logPage, mgr.currentBlock)
	if err != nil {
		return err
	}

	mgr.lastSavedLSN = mgr.latestLSN

	return nil
}

// AppendRecord appends a record to block.
func (mgr *LogManager) AppendRecord(record []byte) (int32, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	boundary, err := mgr.logPage.GetInt32(0)
	if err != nil {
		return 0, err
	}

	bytesNeeded := int32(int32Length + len(record))

	if boundary-bytesNeeded < int32Length {
		if bytesNeeded+int32Length > int32(mgr.fileMgr.BlockSize()) {
			return 0, errors.New("too long record")
		}

		err = mgr.Flush()
		if err != nil {
			return 0, err
		}

		blk, err := mgr.AppendNewBlock()
		if err != nil {
			return 0, err
		}

		mgr.currentBlock = blk
		if err != nil {
			return 0, err
		}

		boundary, err = mgr.logPage.GetInt32(0)
		if err != nil {
			return 0, err
		}
	}

	recordPos := boundary - bytesNeeded
	err = mgr.logPage.SetBytes(int64(recordPos), record)
	if err != nil {
		return 0, err
	}

	err = mgr.logPage.SetInt32(0, recordPos)
	if err != nil {
		return 0, err
	}

	mgr.latestLSN++

	return mgr.latestLSN, nil
}

// AppendNewBlock appends a block to log file and return the appended block.
func (mgr *LogManager) AppendNewBlock() (*Block, error) {
	blk, err := mgr.fileMgr.ExtendFile(mgr.logFileName)
	if err != nil {
		return nil, err
	}

	mgr.logPage.Reset()

	err = mgr.logPage.SetInt32(0, int32(mgr.fileMgr.BlockSize()))
	if err != nil {
		return nil, err
	}

	err = mgr.fileMgr.CopyPageToBlock(mgr.logPage, blk)
	if err != nil {
		return nil, err
	}

	mgr.currentBlock = blk

	return blk, nil
}
