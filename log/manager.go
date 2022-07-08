package log

import (
	"errors"
	"sync"

	"github.com/goropikari/simpledbgo/common"
	"github.com/goropikari/simpledbgo/domain"
)

const (
	boundaryPositionOffset     = 0
	boundaryPositionByteLength = common.Int32Length
	byteStoredByteLength       = common.Int32Length
)

// ManagerConfig is a configuration of log manager.
type ManagerConfig struct {
	LogFileName string
}

type Page struct {
	*domain.Page
}

func NewPage(page *domain.Page) *Page {
	return &Page{
		Page: page,
	}
}

func (p *Page) GetDomainPage() *domain.Page {
	return p.Page
}

func (p *Page) GetBoundaryOffset() (int32, error) {
	return p.GetInt32(boundaryPositionOffset)
}

func (p *Page) setBoundaryOffset(recordPos int32) error {
	return p.SetInt32(boundaryPositionOffset, recordPos)
}

func (p *Page) neededByteLength(recLen int) int32 {
	return int32(byteStoredByteLength + recLen)
}

func (p *Page) canAppend(record []byte) (bool, error) {
	boundary, err := p.GetBoundaryOffset()
	if err != nil {
		return false, err
	}
	bytesNeeded := p.neededByteLength(len(record))

	return boundary-bytesNeeded >= boundaryPositionByteLength, nil
}

func (p *Page) append(record []byte) error {
	boundary, err := p.GetBoundaryOffset()
	if err != nil {
		return err
	}

	bytesNeeded := p.neededByteLength(len(record))

	// 1 record だけで page サイズを超える場合
	if int64(bytesNeeded+boundaryPositionByteLength) > p.Size() {
		return errors.New("too long record")
	}

	recordPos := boundary - bytesNeeded
	if recordPos < boundaryPositionByteLength {
		return errors.New("there is no enough space")
	}

	err = p.SetBytes(int64(recordPos), record)
	if err != nil {
		return err
	}

	err = p.setBoundaryOffset(recordPos)
	if err != nil {
		return err
	}

	return nil
}

// Manager is a log manager.
type Manager struct {
	mu           sync.Mutex
	fileMgr      domain.FileManager
	logFileName  domain.FileName
	currentBlock domain.Block
	logPage      *Page
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
		logPage:      NewPage(page),
		latestLSN:    0,
		lastSavedLSN: 0,
	}, nil
}

func (mgr *Manager) getDomainPage() *domain.Page {
	return mgr.logPage.GetDomainPage()
}

// prepareManager prepares a block and a page for initializing Manager.
// If given file is empty, extend a file by block size.
func prepareManager(fileMgr domain.FileManager, factory *domain.PageFactory, fileName domain.FileName) (domain.Block, *domain.Page, error) {
	page, err := factory.Create()
	if err != nil {
		return domain.Block{}, nil, err
	}

	blklen, err := fileMgr.BlockLength(fileName)
	if err != nil {
		return domain.Block{}, nil, err
	}

	if blklen == 0 {
		blk, err := fileMgr.ExtendFile(fileName)
		if err != nil {
			return domain.Block{}, nil, err
		}

		err = page.SetInt32(boundaryPositionOffset, int32(fileMgr.BlockSize()))
		if err != nil {
			return domain.Block{}, nil, err
		}

		err = fileMgr.CopyPageToBlock(page, blk)
		if err != nil {
			return domain.Block{}, nil, err
		}

		return blk, page, nil
	}

	blknum, err := domain.NewBlockNumber(blklen - 1)
	if err != nil {
		return domain.Block{}, nil, err
	}

	blk := domain.NewBlock(fileName, blknum)

	err = fileMgr.CopyBlockToPage(blk, page)
	if err != nil {
		return domain.Block{}, nil, err
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
	err := mgr.fileMgr.CopyPageToBlock(mgr.getDomainPage(), mgr.currentBlock)
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

	ok, err := mgr.logPage.canAppend(record)
	if err != nil {
		return 0, err
	}
	if !ok {
		if err := mgr.Flush(); err != nil {
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
	}

	if err := mgr.logPage.append(record); err != nil {
		return 0, err
	}

	mgr.latestLSN++

	return mgr.latestLSN, nil
}

// AppendNewBlock appends a block to log file and return the appended block.
func (mgr *Manager) AppendNewBlock() (domain.Block, error) {
	blk, err := mgr.fileMgr.ExtendFile(mgr.logFileName)
	if err != nil {
		return domain.Block{}, err
	}

	mgr.logPage.Reset()

	err = mgr.logPage.SetInt32(boundaryPositionOffset, int32(mgr.fileMgr.BlockSize()))
	if err != nil {
		return domain.Block{}, err
	}

	err = mgr.fileMgr.CopyPageToBlock(mgr.getDomainPage(), blk)
	if err != nil {
		return domain.Block{}, err
	}

	mgr.currentBlock = blk

	return blk, nil
}

// Iterator returns log record iterator.
func (mgr *Manager) Iterator() (domain.LogIterator, error) {
	if err := mgr.Flush(); err != nil {
		return nil, err
	}
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
