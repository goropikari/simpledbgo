package log

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/backend/service"
)

// ErrInvalidArgs is an error that means given args is invalid.
var ErrInvalidArgs = errors.New("arguments is invalid")

// Config is configuration of log manager.
type Config struct {
	logfile core.FileName
}

// NewConfig is constructor of Config.
func NewConfig(logfile core.FileName) Config {
	config := Config{
		logfile: logfile,
	}
	config.SetDefaults()

	return config
}

// SetDefaults sets default value of config.
func (config *Config) SetDefaults() {
	config.logfile = "logfile"
}

// Manager is a log manager of database.
type Manager struct {
	mu           sync.Mutex
	fileMgr      service.FileManager
	currentBlock *core.Block
	page         *core.Page
	latestLSN    int32 // reset when server restart
	lastSavedLSN int32
	config       Config
}

// NewManager is constructor of Manager.
func NewManager(fileMgr service.FileManager, config Config) (*Manager, error) {
	page, err := fileMgr.PreparePage()
	if err != nil {
		return nil, err
	}

	n, err := fileMgr.FileSize(config.logfile)
	if err != nil {
		return nil, err
	}

	// logfile のサイズが 0 だったら block size 分ファイルを作る
	if n == 0 {
		_, _, err := appendNewLogBlock(fileMgr, config.logfile)
		if err != nil {
			return nil, err
		}
	}

	lastBlock, err := fileMgr.LastBlock(config.logfile)
	if err != nil {
		return nil, err
	}

	if err := fileMgr.CopyBlockToPage(lastBlock, page); err != nil {
		return nil, err
	}

	return &Manager{
		mu:           sync.Mutex{},
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
	if err := mgr.fileMgr.CopyPageToBlock(mgr.page, mgr.currentBlock); err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	mgr.lastSavedLSN = mgr.latestLSN

	return nil
}

// FlushByLSN flushes given LSN block.
func (mgr *Manager) FlushByLSN(lsn int32) error {
	if lsn >= mgr.lastSavedLSN {
		return mgr.flush()
	}

	return nil
}

// AppendRecord appends record into the log page.
func (mgr *Manager) AppendRecord(record []byte) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	boundary, err := mgr.page.GetUint32(0)
	if err != nil {
		return fmt.Errorf("failed to get boundary: %w", err)
	}

	recordLength := len(record)
	bytesNeeded := recordLength + core.Uint32Length

	// If there is no enough space, append new block to logfile.
	if int(boundary)-bytesNeeded < core.Uint32Length {
		if err := mgr.flush(); err != nil {
			return err
		}

		page, block, err := appendNewLogBlock(mgr.fileMgr, mgr.config.logfile)
		if err != nil {
			return err
		}

		mgr.page = page
		mgr.currentBlock = block

		boundary, err = mgr.page.GetUint32(0)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	recordPosition := int(boundary) - bytesNeeded

	if err := mgr.page.SetBytes(int64(recordPosition), record); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w", err)
	}

	if err := mgr.page.SetUint32(0, uint32(recordPosition)); err != nil {
		return fmt.Errorf("%w", err)
	}

	mgr.latestLSN++

	return nil
}

// Iterator returns iterator.
func (mgr *Manager) Iterator() (<-chan []byte, error) {
	it, err := iterator(mgr.fileMgr, mgr.currentBlock)
	if err != nil {
		return nil, err
	}

	return it, nil
}

// appendNewLogBlock appends new block to log file.
func appendNewLogBlock(fileMgr service.FileManager, filename core.FileName) (*core.Page, *core.Block, error) {
	page, err := fileMgr.PreparePage()
	if err != nil {
		return nil, nil, err
	}

	blockSize := fileMgr.GetBlockSize()

	if err := page.SetUint32(0, uint32(blockSize)); err != nil {
		return nil, nil, err
	}

	// extend new block
	block, err := fileMgr.AppendBlock(filename)
	if err != nil {
		return nil, nil, err
	}

	if err := fileMgr.CopyPageToBlock(page, block); err != nil && !errors.Is(err, io.EOF) {
		return nil, nil, err
	}

	return page, block, nil
}
