package file

import (
	"io"
	"log"
	stdos "os"
	"path"
	"sync"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/os"
)

const (
	defaultBlockSize = 4096
)

// ManagerConfig is configuration of file manager.
type ManagerConfig struct {
	DBPath    string
	BlockSize int32
	DirectIO  bool
}

// NewManagerConfig constructs a ManagerConfig.
func NewManagerConfig() ManagerConfig {
	c := ManagerConfig{
		DBPath:    path.Join(stdos.Getenv("HOME"), "simpledb"),
		BlockSize: defaultBlockSize,
		DirectIO:  true,
	}

	if path := stdos.Getenv("SIMPLEDB_PATH"); path != "" {
		c.DBPath = path
	}

	return c
}

// Manager is a model of file manager.
type Manager struct {
	mu        sync.Mutex
	explorer  domain.Explorer
	bsf       domain.ByteSliceFactory
	blockSize domain.BlockSize
	dbpath    string
}

// NewManager is a constructor of Manager.
func NewManager(config ManagerConfig) (*Manager, error) {
	var explorer domain.Explorer
	var bsf domain.ByteSliceFactory
	if config.DirectIO {
		explorer = os.NewDirectIOExplorer(config.DBPath) // make server directory
		bsf = bytes.NewDirectByteSliceCreater()
	} else {
		explorer = os.NewNonDirectIOExplorer(config.DBPath) // make server directory
		bsf = bytes.NewByteSliceCreater()
	}

	blkSize, err := domain.NewBlockSize(config.BlockSize)
	if err != nil {
		return nil, errors.Err(err, "create BlockSize")
	}

	return &Manager{
		mu:        sync.Mutex{},
		explorer:  explorer,
		bsf:       bsf,
		blockSize: blkSize,
		dbpath:    config.DBPath,
	}, nil
}

// CopyBlockToPage copies block content to page.
func (mgr *Manager) CopyBlockToPage(blk domain.Block, page *domain.Page) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	page.Reset()

	file, err := mgr.OpenFile(blk.FileName())
	if err != nil {
		return errors.Err(err, "open file")
	}

	_, err = file.Seek(mgr.offset(blk))
	if err != nil {
		return errors.Err(err, "seek")
	}

	// file size が 0 のとき CopyN は EOF を返す。
	// block size 分読んだことにしたいので EOF は無視する。
	_, err = io.CopyN(page, file, int64(mgr.blockSize))
	if err != nil && !errors.Is(err, io.EOF) {
		return errors.Err(err, "copy from file to page")
	}

	_, err = page.Seek(0, io.SeekStart)
	if err != nil {
		return errors.Err(err, "seek")
	}

	return nil
}

// CopyPageToBlock copies page to block.
func (mgr *Manager) CopyPageToBlock(page *domain.Page, block domain.Block) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	file, err := mgr.OpenFile(block.FileName())
	if err != nil {
		return errors.Err(err, "open file")
	}

	if _, err = file.Seek(mgr.offset(block)); err != nil {
		return errors.Err(err, "seek")
	}

	if _, err := file.Write(page.GetData()); err != nil {
		return errors.Err(err, "write")
	}

	if _, err := page.Seek(0, io.SeekStart); err != nil {
		return errors.Err(err, "seek")
	}

	return nil
}

// ExtendFile extends file size by block size and returns last block.
func (mgr *Manager) ExtendFile(filename domain.FileName) (domain.Block, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	blkLen, err := mgr.BlockLength(filename)
	if err != nil {
		return domain.Block{}, errors.Err(err, "take block length")
	}

	numBlk, err := domain.NewBlockNumber(blkLen)
	if err != nil {
		return domain.Block{}, errors.Err(err, "constnruct BlockNumber")
	}

	blk := domain.NewBlock(filename, numBlk)

	file, err := mgr.OpenFile(filename)
	if err != nil {
		return domain.Block{}, errors.Err(err, "open file")
	}

	n, err := file.Size()
	if err != nil {
		return domain.Block{}, errors.Err(err, "take file size")
	}

	_, err = file.Seek(n)
	if err != nil {
		return domain.Block{}, errors.Err(err, "seek")
	}

	bs, err := mgr.bsf.Create(int(mgr.blockSize))
	if err != nil {
		return domain.Block{}, errors.Err(err, "create byte slice")
	}

	_, err = file.Write(bs)
	if err != nil {
		return domain.Block{}, errors.Err(err, "write")
	}

	return blk, nil
}

// BlockLength returns the number of block of the file.
func (mgr *Manager) BlockLength(filename domain.FileName) (int32, error) {
	file, err := mgr.OpenFile(filename)
	if err != nil {
		return 0, errors.Err(err, "open file")
	}

	n, err := file.Size()
	if err != nil {
		return 0, errors.Err(err, "take file size")
	}

	return int32(n) / int32(mgr.blockSize), nil
}

// BlockSize returns block size.
func (mgr *Manager) BlockSize() domain.BlockSize {
	return mgr.blockSize
}

// OpenFile opens a file.
func (mgr *Manager) OpenFile(filename domain.FileName) (*domain.File, error) {
	return mgr.explorer.OpenFile(filename)
}

func (mgr *Manager) offset(blk domain.Block) int64 {
	return int64(mgr.blockSize) * int64(blk.Number())
}

// CreatePage creates a Page.
func (mgr *Manager) CreatePage() (*domain.Page, error) {
	pageFactory := domain.NewPageFactory(mgr.bsf, mgr.blockSize)

	return pageFactory.Create()
}

// IsInit checks whether database is initialized or not.
func (mgr *Manager) IsInit() bool {
	path := mgr.dbpath + "/" + "init"
	_, err := stdos.Stat(path)
	if isNewDatabase := err != nil; isNewDatabase {
		initFile, err := stdos.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		if err := initFile.Close(); err != nil {
			log.Fatal(err)
		}

		return true
	}

	return false
}
