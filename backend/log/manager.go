package log

import (
	"errors"
	"sync"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/meta"
)

// ManagerConfig is a configuration of log manager.
type ManagerConfig struct {
	LogFileName string
}

// Manager is a log manager.
type Manager struct {
	mu           sync.Mutex
	fileMgr      domain.FileManager
	logFileName  domain.FileName
	currentBlock *domain.Block
	logPage      *domain.Page
	pageFactory  *domain.PageFactory
	// Reset when server restarts. Increment when record is appended.
	latestLSN    domain.LSN
	lastSavedLSN domain.LSN
}

// NewManager is a constructor of Manager.
func NewManager(fileMgr domain.FileManager, pageFactory *domain.PageFactory, config ManagerConfig) (*Manager, error) {
	logFileName, err := domain.NewFileName(config.LogFileName)
	if err != nil {
		return nil, err
	}

	block, page, err := prepareManager(fileMgr, pageFactory, logFileName)
	if err != nil {
		return nil, err
	}

	return &Manager{
		mu:           sync.Mutex{},
		fileMgr:      fileMgr,
		logFileName:  logFileName,
		pageFactory:  pageFactory,
		currentBlock: block,
		logPage:      page,
		latestLSN:    0,
		lastSavedLSN: 0,
	}, nil
}

// prepareManager prepares a block and a page for initializing Manager.
// If given file is empty, extend a file by block size.
func prepareManager(fileMgr domain.FileManager, factory *domain.PageFactory, fileName domain.FileName) (*domain.Block, *domain.Page, error) {
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
func (mgr *Manager) FlushLSN(lsn domain.LSN) error {
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
func (mgr *Manager) AppendRecord(record []byte) (domain.LSN, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	boundary, err := mgr.logPage.GetInt32(0)
	if err != nil {
		return 0, err
	}

	bytesNeeded := int32(meta.Int32Length + len(record))

	if bytesNeeded+meta.Int32Length > int32(mgr.fileMgr.BlockSize()) {
		return 0, errors.New("too long record")
	}

	if boundary-bytesNeeded < meta.Int32Length {
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

// Iterator returns log record iterator.
func (mgr *Manager) Iterator() (domain.LogIterator, error) {
	page, err := mgr.pageFactory.Create()
	if err != nil {
		return nil, err
	}

	return NewIterator(mgr.fileMgr, mgr.currentBlock, page)
}

// LogFileName returns log file name.
func (mgr *Manager) LogFileName() domain.FileName {
	return mgr.logFileName
}
