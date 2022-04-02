package log

import (
	"errors"
	"sync"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/backend/file"
)

const int32Length = 4

// ManagerConfig is a configuration of log manager.
type ManagerConfig struct {
	LogFileName string
}

// Manager is a log manager.
type Manager struct {
	mu           sync.Mutex
	fileMgr      *file.Manager
	logFileName  domain.FileName
	currentBlock *domain.Block
	logPage      *domain.Page
	// Reset when server restarts. Increment when record is appended.
	latestLSN    int32
	lastSavedLSN int32
}

// NewManager is a constructor of Manager.
func NewManager(fileMgr *file.Manager, block *domain.Block, page *domain.Page, config ManagerConfig) (*Manager, error) {
	logFileName, err := domain.NewFileName(config.LogFileName)
	if err != nil {
		return nil, err
	}

	return &Manager{
		mu:           sync.Mutex{},
		fileMgr:      fileMgr,
		logFileName:  logFileName,
		currentBlock: block,
		logPage:      page,
		latestLSN:    0,
		lastSavedLSN: 0,
	}, nil
}

// PrepareManager prepares a block and a page for initializing Manager.
// If given file is empty, extend a file by block size.
func PrepareManager(fileMgr *file.Manager, factory *domain.PageFactory, fileName domain.FileName) (*domain.Block, *domain.Page, error) {
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

	blknum, err := domain.NewBlockNumber(blklen - 1)
	if err != nil {
		return nil, nil, err
	}

	blk := domain.NewBlock(fileName, fileMgr.BlockSize(), blknum)

	err = fileMgr.CopyBlockToPage(blk, page)
	if err != nil {
		return nil, nil, err
	}

	return blk, page, nil
}

// FlushLSN flushes by lsn.
func (mgr *Manager) FlushLSN(lsn int32) error {
	if lsn >= mgr.lastSavedLSN {
		return mgr.Flush()
	}

	return nil
}

// Flush flushes the log page.
func (mgr *Manager) Flush() error {
	err := mgr.fileMgr.CopyPageToBlock(mgr.logPage, mgr.currentBlock)
	if err != nil {
		return err
	}

	mgr.lastSavedLSN = mgr.latestLSN

	return nil
}

// AppendRecord appends a record to block.
func (mgr *Manager) AppendRecord(record []byte) (int32, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	boundary, err := mgr.logPage.GetInt32(0)
	if err != nil {
		return 0, err
	}

	bytesNeeded := int32(int32Length + len(record))

	if bytesNeeded+int32Length > int32(mgr.fileMgr.BlockSize()) {
		return 0, errors.New("too long record")
	}

	if boundary-bytesNeeded < int32Length {
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
func (mgr *Manager) AppendNewBlock() (*domain.Block, error) {
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
