package domain

import (
	"sync"
)

type LogManagerConfig struct {
	LogFileName string
}

type LogManager struct {
	mu           sync.Mutex
	fileMgr      *FileManager
	logFileName  FileName
	currentBlock *Block
	logPage      *Page
	lastestLSN   int32
	lastSavedLSN int32
}

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
		lastestLSN:   0,
		lastSavedLSN: 0,
	}, nil
}

func (mgr *LogManager) FlushLSN(lsn int32) error {
	return nil
}

func (mgr *LogManager) AppendRecord(record []byte) error {
	return nil
}

func (mgr *LogManager) AppendNewBlock() (*Block, error) {
	block, err := mgr.fileMgr.ExtendFile(mgr.logFileName)
	if err != nil {
		return nil, err
	}

	mgr.logPage.Reset()

	err = mgr.logPage.SetInt32(0, int32(mgr.fileMgr.blockSize))
	if err != nil {
		return nil, err
	}

	err = mgr.fileMgr.CopyPageToBlock(mgr.logPage, block)
	if err != nil {
		return nil, err
	}

	mgr.currentBlock = block

	return block, nil
}

func (mgr *LogManager) FlushPage() error {
	return nil
}
