package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goropikari/simpledb_go/core"
	"github.com/goropikari/simpledb_go/directio"
)

// Config is configuration of Manager.
type Config struct {
	dbDir      string
	blockSize  int // for direct io, blockSize must be multiple of 4096
	isDirectIO bool
}

// NewConfig is constructor of Config
func NewConfig(dbDir string, blockSize int, isDirectIO bool) (Config, error) {
	if isDirectIO && blockSize%directio.BlockSize != 0 {
		return Config{}, directio.InvalidBlockSize
	}

	abspath, err := filepath.Abs(dbDir)
	if err != nil {
		return Config{}, err
	}

	return Config{
		dbDir:      abspath,
		blockSize:  blockSize,
		isDirectIO: isDirectIO,
	}, nil

}

// Manager manages files.
type Manager struct {
	mu        sync.Mutex
	config    Config
	openFiles map[core.FileName]*os.File
}

// NewManager is constructor of Manager.
func NewManager(config Config) (*Manager, error) {
	if err := os.MkdirAll(config.dbDir, os.ModePerm); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(config.dbDir)
	if err != nil {
		return nil, err
	}

	// remove temporary files.
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "temp") {
			if err := os.Remove(filepath.Join(config.dbDir, file.Name())); err != nil {
				return nil, err
			}
		}
	}

	return &Manager{
		config:    config,
		openFiles: make(map[core.FileName]*os.File, 0),
	}, nil
}

// CopyBlockToPage copies block to page.
func (fileMgr *Manager) CopyBlockToPage(block *Block, page *Page) error {
	if fileMgr == nil {
		return core.NilReceiverError
	}

	fileMgr.mu.Lock()
	defer fileMgr.mu.Unlock()

	f, err := fileMgr.openFile(block.GetFileName())
	if err != nil {
		return err
	}

	if _, err = f.Seek(int64(fileMgr.config.blockSize*int(block.GetBlockNumber())), io.SeekStart); err != nil {
		return err
	}

	if _, err = io.CopyN(page, f, int64(fileMgr.config.blockSize)); err != nil {
		return err
	}

	return nil
}

// CopyPageToBlock copies page to block.
func (fileMgr *Manager) CopyPageToBlock(page *Page, block *Block) error {
	if fileMgr == nil {
		return core.NilReceiverError
	}

	fileMgr.mu.Lock()
	defer fileMgr.mu.Unlock()

	f, err := fileMgr.openFile(block.GetFileName())
	if err != nil {
		return err
	}

	if _, err = f.Seek(int64(fileMgr.config.blockSize*int(block.GetBlockNumber())), io.SeekStart); err != nil {
		return err
	}

	if _, err := f.Write(page.GetFullBytes()); err != nil {
		return err
	}

	return nil
}

// AppendBlock appends block to given filename.
func (fileMgr *Manager) AppendBlock(filename core.FileName) (*Block, error) {
	if fileMgr == nil {
		return nil, core.NilReceiverError
	}

	fileMgr.mu.Lock()
	defer fileMgr.mu.Unlock()

	f, err := fileMgr.openFile(filename)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if _, err = f.Seek(stat.Size(), io.SeekStart); err != nil {
		return nil, err
	}

	blockNum := core.BlockNumber(stat.Size() / int64(fileMgr.config.blockSize))
	block := NewBlock(filename, blockNum)
	blk, err := fileMgr.prepareBytes()
	if err != nil {
		return nil, err
	}

	if _, err = f.Write(blk); err != nil {
		return nil, err
	}

	return block, nil
}

func (fileMgr *Manager) openFile(filename core.FileName) (f *os.File, err error) {
	if fileMgr == nil {
		return nil, core.NilReceiverError
	}

	if v, ok := fileMgr.openFiles[filename]; ok {
		return v, nil
	}

	// open file. If there is no such file, create new file.
	path := filepath.Join(string(fileMgr.config.dbDir), string(filename))
	if fileMgr.config.isDirectIO {
		f, err = directio.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else {
		f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	fileMgr.openFiles[filename] = f

	return f, nil
}

func (fileMgr *Manager) prepareBytes() (blk []byte, err error) {
	if fileMgr.config.isDirectIO {
		blk, err = directio.AlignedBlock(int(fileMgr.config.blockSize))
		if err != nil {
			return nil, err
		}
	} else {
		blk = make([]byte, fileMgr.config.blockSize)
	}

	return
}

// CloseFile closes a file.
func (fileMgr *Manager) CloseFile(filename core.FileName) error {
	if fileMgr == nil {
		return core.NilReceiverError
	}

	if f, ok := fileMgr.openFiles[filename]; ok {
		delete(fileMgr.openFiles, filename)
		if err := f.Close(); err != nil {
			return err
		}
		return nil
	}

	return errors.New("no such file")
}
