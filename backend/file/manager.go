package file

import (
	"errors"
	"fmt"
	"io"
	goos "os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/infra"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/lib/os"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// Explorer is an interface of file explorer.
type Explorer interface {
	MkdirAll(path string) error
	ReadDir(name string) ([]goos.DirEntry, error)
	Remove(dir string, file string) error
	OpenFile(path string) (*os.File, error)
}

var (
	// ErrNoSuchFile is an error that means specified file does not exist.
	ErrNoSuchFile = errors.New("no such file")

	// ErrInvalidArgs is an error that means given args is invalid.
	ErrInvalidArgs = errors.New("arguments is invalid")

	// ErrInvalidConfig is an error that means given config is invalid.
	ErrInvalidConfig = errors.New("config is invalid")
)

// Manager manages files.
type Manager struct {
	mu        sync.Mutex
	config    infra.Config
	explorer  Explorer
	openFiles map[core.FileName]*os.File
}

// NewManager is constructor of Manager.
func NewManager(exp Explorer, config infra.Config) (*Manager, error) {
	config.SetDefaults()

	if err := exp.MkdirAll(config.DBPath); err != nil {
		return nil, err
	}

	if err := deleteTempFiles(exp, config.DBPath); err != nil {
		return nil, err
	}

	return &Manager{
		mu:        sync.Mutex{},
		config:    config,
		explorer:  exp,
		openFiles: make(map[core.FileName]*os.File, 0),
	}, nil
}

func deleteTempFiles(exp Explorer, dbPath string) error {
	files, err := exp.ReadDir(dbPath)
	if err != nil {
		return err
	}

	// remove temporary files.
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "temp") {
			if err := exp.Remove(dbPath, file.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetBlockSize returns block size.
func (mgr *Manager) GetBlockSize() int {
	return mgr.config.BlockSize
}

// CopyBlockToPage copies block to page.
func (mgr *Manager) CopyBlockToPage(block *core.Block, page *core.Page) error {
	if block == nil {
		return ErrInvalidArgs
	}

	if page == nil {
		return ErrInvalidArgs
	}

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	file, err := mgr.openFile(block.GetFileName())
	if err != nil {
		return err
	}

	if _, err := page.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	// seek でファイルサイズ以上の位置が指定されていた場合、io.CopyN しても 1 byte も読み込まれず
	// page に変化がない.
	// 実際は x00 を blocksize 分読み込んだということにしたいので、page を 0 reset しておく
	page.Reset()

	seekPos := int64(mgr.config.BlockSize * int(block.GetBlockNumber()))

	if _, err = file.Seek(seekPos, io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err = io.CopyN(page, file, int64(mgr.config.BlockSize)); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w", err)
	}

	if _, err := page.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// CopyPageToBlock copies page to block.
func (mgr *Manager) CopyPageToBlock(page *core.Page, block *core.Block) error {
	if block == nil {
		return ErrInvalidArgs
	}

	if page == nil {
		return ErrInvalidArgs
	}

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	file, err := mgr.openFile(block.GetFileName())
	if err != nil {
		return err
	}

	if _, err = file.Seek(int64(mgr.config.BlockSize*int(block.GetBlockNumber())), io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := file.Write(page.GetBufferBytes()); err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := page.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// AppendBlock appends block to given filename.
func (mgr *Manager) AppendBlock(filename core.FileName) (*core.Block, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	file, err := mgr.openFile(filename)
	if err != nil {
		return nil, err
	}

	numBlock, err := mgr.numBlock(file)
	if err != nil {
		return nil, err
	}

	appendBlockNum, err := core.NewBlockNumber(numBlock)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	block := core.NewBlock(filename, appendBlockNum)

	// extend file size
	offset := int64(numBlock * mgr.config.BlockSize)
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	page, err := mgr.PreparePage()
	if err != nil {
		return nil, err
	}

	if _, err = file.Write(page.GetBufferBytes()); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return block, nil
}

// numBlock returns the number of blocks of given file.
func (mgr *Manager) numBlock(file *os.File) (int, error) {
	fileSize, err := file.Size()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return int(fileSize) / mgr.config.BlockSize, nil
}

// LastBlock returns last block of given file.
func (mgr *Manager) LastBlock(filename core.FileName) (*core.Block, error) {
	f, err := mgr.openFile(filename)
	if err != nil {
		return nil, err
	}

	nb, err := mgr.numBlock(f)
	if err != nil {
		return nil, err
	}

	lastBlockNumber, err := core.NewBlockNumber(nb - 1)
	if err != nil {
		return nil, err
	}

	return core.NewBlock(filename, lastBlockNumber), nil
}

// openFile opens file as given filename.
// If there is no such file, create new file.
func (mgr *Manager) openFile(filename core.FileName) (*os.File, error) {
	if v, ok := mgr.openFiles[filename]; ok {
		return v, nil
	}

	path := filepath.Join(mgr.config.DBPath, string(filename))

	// open file. If there is no such file, create new file.
	f, err := mgr.explorer.OpenFile(path)
	if err != nil {
		return nil, err
	}

	mgr.openFiles[filename] = f

	return f, nil
}

// CloseFile closes a file.
func (mgr *Manager) closeFile(filename core.FileName) error {
	if f, ok := mgr.openFiles[filename]; ok {
		delete(mgr.openFiles, filename)

		if err := f.Close(); err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}

	return ErrNoSuchFile
}

func (mgr *Manager) FileSize(filename core.FileName) (int64, error) {
	file, err := mgr.openFile(filename)
	if err != nil {
		return 0, err
	}

	return file.Size()
}

// PreparePage prepares a page.
// If file manager's config specifies direct IO support, this returns page
// satisfying direct IO constraints.
func (mgr *Manager) PreparePage() (*core.Page, error) {
	if mgr.config.IsDirectIO {
		bb, err := bytes.NewDirectBuffer(mgr.config.BlockSize)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return core.NewPage(bb), nil
	}

	bb := bytes.NewBuffer(mgr.config.BlockSize)

	return core.NewPage(bb), nil
}
