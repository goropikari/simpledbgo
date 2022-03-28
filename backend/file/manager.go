package file

import (
	"errors"
	"fmt"
	"io"
	"log"
	goos "os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goropikari/simpledb_go/backend/core"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/lib/directio"
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

// Config is configuration of Manager.
type Config struct {
	dbDir      string
	blockSize  int // for direct io, blockSize must be multiple of 4096
	isDirectIO bool
}

// NewConfig is constructor of Config.
func NewConfig(dbDir string, blockSize int, isDirectIO bool) Config {
	if isDirectIO && blockSize%directio.BlockSize != 0 {
		log.Fatal(directio.ErrInvalidBlockSize)
	}

	abspath, err := filepath.Abs(dbDir)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{
		dbDir:      abspath,
		blockSize:  blockSize,
		isDirectIO: isDirectIO,
	}

	return config
}

// SetDefaults sets defalut value of config.
func (config *Config) SetDefaults() {
	if config.dbDir == "" {
		abspath, _ := filepath.Abs("simpledb")
		config.dbDir = abspath
	}

	if config.blockSize == 0 {
		config.blockSize = directio.BlockSize
	}
}

// Manager manages files.
type Manager struct {
	mu        sync.Mutex
	config    Config
	explorer  Explorer
	openFiles map[core.FileName]*os.File
}

// NewManager is constructor of Manager.
func NewManager(exp Explorer, config Config) *Manager {
	config.SetDefaults()

	if err := exp.MkdirAll(config.dbDir); err != nil {
		log.Fatal(err)
	}

	if err := deleteTempFiles(exp, config.dbDir); err != nil {
		log.Fatal(err)
	}

	return &Manager{
		mu:        sync.Mutex{},
		config:    config,
		explorer:  exp,
		openFiles: make(map[core.FileName]*os.File, 0),
	}
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
	return mgr.config.blockSize
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

	seekPos := int64(mgr.config.blockSize * int(block.GetBlockNumber()))

	if _, err = file.Seek(seekPos, io.SeekStart); err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err = io.CopyN(page, file, int64(mgr.config.blockSize)); err != nil && !errors.Is(err, io.EOF) {
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

	if _, err = file.Seek(int64(mgr.config.blockSize*int(block.GetBlockNumber())), io.SeekStart); err != nil {
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
	offset := int64(numBlock * mgr.config.blockSize)
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	buf := mgr.prepareBytes()

	if _, err = file.Write(buf); err != nil {
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

	return int(fileSize) / mgr.config.blockSize, nil
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

	path := filepath.Join(mgr.config.dbDir, string(filename))

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

// prepareBytes prepares byte slice.
func (mgr *Manager) prepareBytes() []byte {
	if mgr.config.isDirectIO {
		buf, err := directio.AlignedBlock(mgr.config.blockSize)
		if err != nil {
			log.Fatal(err)
		}

		return buf
	}

	return make([]byte, mgr.config.blockSize)
}

// PreparePage prepares a page.
// If file manager's config specifies direct IO support, this returns page
// satisfying direct IO constraints.
func (mgr *Manager) PreparePage() (*core.Page, error) {
	if mgr.config.isDirectIO {
		bb, err := bytes.NewDirectBuffer(mgr.config.blockSize)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return core.NewPage(bb), nil
	}

	bb := bytes.NewBuffer(mgr.config.blockSize)

	return core.NewPage(bb), nil
}