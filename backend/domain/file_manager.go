package domain

import (
	"fmt"
	"io"
	"sync"
)

// FileManagerConfig is configuration of file manager.
type FileManagerConfig struct {
	BlockSize int32
}

// FileManager is a model of file manager.
type FileManager struct {
	mu        sync.Mutex
	explorer  Explorer
	bsf       ByteSliceFactory
	blockSize BlockSize
}

// NewFileManager is a constructor of FileManager.
func NewFileManager(explorer Explorer, bsf ByteSliceFactory, config FileManagerConfig) (*FileManager, error) {
	blkSize, err := NewBlockSize(config.BlockSize)
	if err != nil {
		return nil, err
	}

	return &FileManager{
		mu:        sync.Mutex{},
		explorer:  explorer,
		bsf:       bsf,
		blockSize: blkSize,
	}, nil
}

// CopyBlockToPage copies block content to page.
func (mgr *FileManager) CopyBlockToPage(blk *Block, page *Page) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	page.Reset()

	file, err := mgr.OpenFile(blk.FileName())
	if err != nil {
		return err
	}

	_, err = file.Seek(blk.Offset())
	if err != nil {
		return err
	}

	_, err = io.CopyN(page, file, int64(blk.Size()))
	if err != nil {
		return err
	}

	_, err = page.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}

// CopyPageToBlock copies page to block.
func (mgr *FileManager) CopyPageToBlock(page *Page, block *Block) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	file, err := mgr.OpenFile(block.FileName())
	if err != nil {
		return err
	}

	if _, err = file.Seek(block.Offset()); err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := file.Write(page.GetData()); err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := page.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// ExtendFile extends file size by block size and returns last block.
func (mgr *FileManager) ExtendFile(filename FileName) (*Block, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	blkLen, err := mgr.BlockLength(filename)
	if err != nil {
		return nil, err
	}

	numBlk, err := NewBlockNumber(blkLen)
	if err != nil {
		return nil, err
	}

	blk := NewBlock(filename, mgr.blockSize, numBlk)

	file, err := mgr.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	n, err := file.Size()
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(n)
	if err != nil {
		return nil, err
	}

	bs, err := mgr.bsf.Create(int(mgr.blockSize))
	if err != nil {
		return nil, err
	}

	_, err = file.Write(bs)
	if err != nil {
		return nil, err
	}

	return blk, nil
}

// BlockLength returns the number of block of the file.
func (mgr *FileManager) BlockLength(filename FileName) (int32, error) {
	file, err := mgr.OpenFile(filename)
	if err != nil {
		return 0, err
	}

	n, err := file.Size()
	if err != nil {
		return 0, err
	}

	return int32(n) / int32(mgr.blockSize), nil
}

func (mgr *FileManager) BlockSize() BlockSize {
	return mgr.blockSize
}

// OpenFile opens a file.
func (mgr *FileManager) OpenFile(filename FileName) (*File, error) {
	return mgr.explorer.OpenFile(filename)
}
