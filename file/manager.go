package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goropikari/simpledb_go/bytes"
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
		return Config{}, directio.InvalidBlockSizeError
	}

	if blockSize <= 0 {
		return Config{}, errors.New("block size must be positive")
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
	if err := validateConfig(config); err != nil {
		return nil, err
	}

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

func validateConfig(config Config) error {
	if config.dbDir == "" {
		return errors.New("database directory must be specified")
	}
	if config.blockSize <= 0 {
		return errors.New("block size must be positive")
	}

	return nil
}

// GetBlockSize returns block size.
func (mgr *Manager) GetBlockSize() (int, error) {
	if mgr == nil {
		return 0, core.NilReceiverError
	}

	return mgr.config.blockSize, nil
}

// CopyBlockToPage copies block to page.
func (mgr *Manager) CopyBlockToPage(block *Block, page *Page) error {
	if mgr == nil {
		return core.NilReceiverError
	}
	if block == nil {
		return errors.New("block must not be nil")
	}
	if page == nil {
		return errors.New("page must not be nil")
	}

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	f, err := mgr.openFile(block.GetFileName())
	if err != nil {
		return err
	}

	if _, err := page.bb.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// seek でファイルサイズ以上の位置が指定されていた場合、io.CopyN しても 1 byte も読み込まれず
	// page に変化がない.
	// 実際は 0 を blocksize 分読み込んだということにしたいので、page を 0 reset しておく
	page.Reset()
	seekPos := int64(mgr.config.blockSize * int(block.GetBlockNumber()))
	if _, err = f.Seek(seekPos, io.SeekStart); err != nil {
		return err
	}

	if _, err = io.CopyN(page, f, int64(mgr.config.blockSize)); err != nil && err != io.EOF {
		return err
	}

	if _, err := page.bb.Seek(0, io.SeekStart); err != nil {
		return err
	}

	return nil
}

// CopyPageToBlock copies page to block.
func (mgr *Manager) CopyPageToBlock(page *Page, block *Block) error {
	if mgr == nil {
		return core.NilReceiverError
	}
	if block == nil {
		return errors.New("block must not be nil")
	}
	if page == nil {
		return errors.New("page must not be nil")
	}

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	f, err := mgr.openFile(block.GetFileName())
	if err != nil {
		return err
	}

	if _, err = f.Seek(int64(mgr.config.blockSize*int(block.GetBlockNumber())), io.SeekStart); err != nil {
		return err
	}

	if _, err := f.Write(page.GetFullBytes()); err != nil {
		return err
	}
	if _, err := page.bb.Seek(0, io.SeekStart); err != nil {
		return err
	}

	return nil
}

// AppendBlock appends block to given filename.
func (mgr *Manager) AppendBlock(filename core.FileName) (*Block, error) {
	if mgr == nil {
		return nil, core.NilReceiverError
	}

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	f, err := mgr.openFile(filename)
	if err != nil {
		return nil, err
	}

	numBlock, err := mgr.numBlock(f)
	if err != nil {
		return nil, err
	}
	appendBlockNum, err := core.NewBlockNumber(numBlock)
	if err != nil {
		return nil, err
	}
	block := NewBlock(filename, appendBlockNum)

	buf, err := mgr.prepareBytes()
	if err != nil {
		return nil, err
	}

	// extend file size
	if _, err := f.Seek(int64(numBlock*mgr.config.blockSize), io.SeekStart); err != nil {
		return nil, err
	}
	if _, err = f.Write(buf); err != nil {
		return nil, err
	}

	return block, nil
}

func (mgr *Manager) numBlock(file *os.File) (int, error) {
	fileSize, err := core.FileSize(file)
	if err != nil {
		return 0, err
	}

	return int(fileSize) / mgr.config.blockSize, nil
}

// LastBlock returns last block of given file.
func (mgr *Manager) LastBlock(filename core.FileName) (*Block, error) {
	if mgr == nil {
		return nil, core.NilReceiverError
	}

	f, err := mgr.openFile(filename)
	if err != nil {
		return nil, err
	}

	lastBlockNum, err := mgr.lastBlockNumber(f)
	if err != nil {
		return nil, err
	}

	block := NewBlock(filename, lastBlockNum)

	return block, nil
}

// lastBlockNumber returns last block number of given file.
func (mgr *Manager) lastBlockNumber(f *os.File) (core.BlockNumber, error) {
	fileSize, err := core.FileSize(f)
	if err != nil {
		return 0, err
	}
	if fileSize == 0 {
		return core.NewBlockNumber(0)
	}

	blockNum, err := core.NewBlockNumber(int(fileSize/int64(mgr.config.blockSize)) - 1)
	if err != nil {
		return 0, err
	}

	return blockNum, nil
}

// openFile opens file as given filename.
// If there is no such file, create new file.
func (mgr *Manager) openFile(filename core.FileName) (f *os.File, err error) {
	if mgr == nil {
		return nil, core.NilReceiverError
	}

	if v, ok := mgr.openFiles[filename]; ok {
		return v, nil
	}

	// open file. If there is no such file, create new file.
	path := filepath.Join(string(mgr.config.dbDir), string(filename))
	flag := os.O_RDWR | os.O_CREATE
	if mgr.config.isDirectIO {
		f, err = directio.OpenFile(path, flag, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else {
		f, err = os.OpenFile(path, flag, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	mgr.openFiles[filename] = f

	return f, nil
}

// CloseFile closes a file.
func (mgr *Manager) CloseFile(filename core.FileName) error {
	if mgr == nil {
		return core.NilReceiverError
	}

	if f, ok := mgr.openFiles[filename]; ok {
		delete(mgr.openFiles, filename)
		if err := f.Close(); err != nil {
			return err
		}
		return nil
	}

	return errors.New("there is no such file")
}

// prepareBytes prepares byte slice.
func (mgr *Manager) prepareBytes() (buf []byte, err error) {
	if mgr.config.isDirectIO {
		buf, err = directio.AlignedBlock(int(mgr.config.blockSize))
		if err != nil {
			return nil, err
		}
	} else {
		buf = make([]byte, mgr.config.blockSize)
	}

	return
}

// PreparePage prepares a page.
// If file manager's config specifies direct IO support, this returns page
// satisfying direct IO constraints.
func (mgr *Manager) PreparePage() (*Page, error) {
	if mgr == nil {
		return nil, core.NilReceiverError
	}

	if mgr.config.isDirectIO {
		bb, err := bytes.NewDirectBuffer(mgr.config.blockSize)
		if err != nil {
			return nil, err
		}

		return NewPage(bb), nil
	}

	bb, err := bytes.NewBuffer(mgr.config.blockSize)
	if err != nil {
		return nil, err
	}

	return NewPage(bb), nil
}
