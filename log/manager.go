package log

import (
	"errors"
	"io"
	"sync"

	"github.com/goropikari/simpledb_go/core"
	"github.com/goropikari/simpledb_go/file"
)

// Config is configuration of log manager.
type Config struct {
	logfile core.FileName
}

// NewConfig is constructor of Config.
func NewConfig(logfile core.FileName) Config {
	return Config{
		logfile: logfile,
	}
}

// Manager is a log manager of database.
type Manager struct {
	mu           sync.Mutex
	fileMgr      *file.Manager
	currentBlock *file.Block
	page         *file.Page
	latestLSN    int // reset when server restart
	lastSavedLSN int
	config       Config
}

// NewManager is constructor of Manager.
func NewManager(fileMgr *file.Manager, config Config) (*Manager, error) {
	if fileMgr == nil {
		return nil, errors.New("fileMgr must not be nil")
	}

	page, err := fileMgr.PreparePage()
	if err != nil {
		return nil, err
	}
	lastBlock, err := fileMgr.LastBlock(config.logfile)
	if err != nil {
		return nil, err
	}

	if err := fileMgr.CopyBlockToPage(lastBlock, page); err != nil {
		return nil, err
	}

	boundary, err := page.GetUint32(0)
	if err != nil {
		return nil, err
	}
	if boundary == 0 {
		blockSize := fileMgr.GetBlockSize()
		if err != nil {
			return nil, err
		}
		if err = page.SetUint32(0, uint32(blockSize)); err != nil {
			return nil, err
		}
	}

	return &Manager{
		fileMgr:      fileMgr,
		currentBlock: lastBlock,
		page:         page,
		latestLSN:    0,
		lastSavedLSN: 0,
		config:       config,
	}, nil
}

// flush flushes page into current block.
func (mgr *Manager) flush() error {
	if mgr == nil {
		return nil
	}

	if err := mgr.fileMgr.CopyPageToBlock(mgr.page, mgr.currentBlock); err != nil {
		return err
	}

	mgr.lastSavedLSN = mgr.latestLSN

	return nil
}

// FlushByLSN flushes given LSN block.
func (mgr *Manager) FlushByLSN(lsn int) error {
	if lsn >= mgr.lastSavedLSN {
		return mgr.flush()
	}

	return nil
}

// AppendRecord appends record into the log page.
func (mgr *Manager) AppendRecord(record []byte) error {
	if mgr == nil {
		return nil
	}

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	boundary, err := mgr.page.GetUint32(0)
	if err != nil {
		return err
	}
	recordLength := len(record)
	bytesNeeded := recordLength + core.Uint32Length

	if int(boundary)-bytesNeeded < core.Uint32Length {
		mgr.flush()
		if err := mgr.appendNewLogBlock(); err != nil {
			return err
		}
		boundary, err = mgr.page.GetUint32(0)
		if err != nil {
			return err
		}
	}
	recordPosition := int(boundary) - bytesNeeded

	if err := mgr.page.SetBytes(int64(recordPosition), record); err != nil && err != io.EOF {
		return err
	}

	if err := mgr.page.SetUint32(0, uint32(recordPosition)); err != nil {
		return err
	}
	mgr.latestLSN++

	return nil
}

// Iterator returns iterator.
func (mgr *Manager) Iterator() (<-chan []byte, error) {
	if mgr == nil {
		return nil, nil
	}

	it, err := iterator(mgr.fileMgr, mgr.currentBlock)
	if err != nil {
		return nil, err
	}

	return it, nil
}

// appendNewLogBlock appends new block to log file.
// This initializes page and append new block to log file.
func (mgr *Manager) appendNewLogBlock() error {
	if mgr == nil {
		return nil
	}

	// flush page into current block
	if err := mgr.flush(); err != nil {
		return err
	}

	// initialize page for log
	blockSize := mgr.fileMgr.GetBlockSize()
	if err := mgr.page.SetUint32(0, uint32(blockSize)); err != nil {
		return err
	}

	// extend new block
	block, err := mgr.fileMgr.AppendBlock(mgr.config.logfile)
	if err != nil {
		return err
	}
	mgr.currentBlock = block

	if err := mgr.fileMgr.CopyPageToBlock(mgr.page, block); err != nil && err != io.EOF {
		return err
	}

	return nil
}
