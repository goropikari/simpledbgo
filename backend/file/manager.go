package file

import (
	"io"
	"sync"

	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/pkg/errors"
)

// ManagerConfig is configuration of file manager.
type ManagerConfig struct {
	BlockSize int32
}

// Manager is a model of file manager.
type Manager struct {
	mu        sync.Mutex
	explorer  domain.Explorer
	bsf       domain.ByteSliceFactory
	blockSize domain.BlockSize
}

// NewManager is a constructor of Manager.
func NewManager(explorer domain.Explorer, bsf domain.ByteSliceFactory, config ManagerConfig) (*Manager, error) {
	blkSize, err := domain.NewBlockSize(config.BlockSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create BlockSize")
	}

	return &Manager{
		mu:        sync.Mutex{},
		explorer:  explorer,
		bsf:       bsf,
		blockSize: blkSize,
	}, nil
}

// CopyBlockToPage copies block content to page.
func (mgr *Manager) CopyBlockToPage(blk domain.Block, page *domain.Page) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	page.Reset()

	file, err := mgr.OpenFile(blk.FileName())
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	_, err = file.Seek(blk.Offset())
	if err != nil {
		return errors.Wrap(err, "failed to seek")
	}

	// file size が 0 のとき CopyN は EOF を返す。
	// block size 分読んだことにしたいので EOF は無視する。
	_, err = io.CopyN(page, file, int64(blk.Size()))
	if err != nil && !errors.Is(err, io.EOF) {
		return errors.Wrap(err, "failed to copy from file to page")
	}

	_, err = page.Seek(0, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "failed to seek")
	}

	return nil
}

// CopyPageToBlock copies page to block.
func (mgr *Manager) CopyPageToBlock(page *domain.Page, block domain.Block) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	file, err := mgr.OpenFile(block.FileName())
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	if _, err = file.Seek(block.Offset()); err != nil {
		return errors.Wrap(err, "failed to seek")
	}

	if _, err := file.Write(page.GetData()); err != nil {
		return errors.Wrap(err, "failed to write")
	}

	if _, err := page.Seek(0, io.SeekStart); err != nil {
		return errors.Wrap(err, "failed to seek")
	}

	return nil
}

// ExtendFile extends file size by block size and returns last block.
func (mgr *Manager) ExtendFile(filename domain.FileName) (domain.Block, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	blkLen, err := mgr.BlockLength(filename)
	if err != nil {
		return domain.NewZeroBlock(), errors.Wrap(err, "failed to take block length")
	}

	numBlk, err := domain.NewBlockNumber(blkLen)
	if err != nil {
		return domain.NewZeroBlock(), errors.Wrap(err, "failed to constnruct BlockNumber")
	}

	blk := domain.NewBlock(filename, mgr.blockSize, numBlk)

	file, err := mgr.OpenFile(filename)
	if err != nil {
		return domain.NewZeroBlock(), errors.Wrap(err, "failed to open file")
	}

	n, err := file.Size()
	if err != nil {
		return domain.NewZeroBlock(), errors.Wrap(err, "failed to take file size")
	}

	_, err = file.Seek(n)
	if err != nil {
		return domain.NewZeroBlock(), errors.Wrap(err, "failed to seek")
	}

	bs, err := mgr.bsf.Create(int(mgr.blockSize))
	if err != nil {
		return domain.NewZeroBlock(), errors.Wrap(err, "failed to create byte slice")
	}

	_, err = file.Write(bs)
	if err != nil {
		return domain.NewZeroBlock(), errors.Wrap(err, "failed to write")
	}

	return blk, nil
}

// BlockLength returns the number of block of the file.
func (mgr *Manager) BlockLength(filename domain.FileName) (int32, error) {
	file, err := mgr.OpenFile(filename)
	if err != nil {
		return 0, errors.Wrap(err, "failed to open file")
	}

	n, err := file.Size()
	if err != nil {
		return 0, errors.Wrap(err, "failed to take file size")
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
