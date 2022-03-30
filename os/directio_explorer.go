package os

import (
	"os"
	"path/filepath"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/lib/directio"
)

// DirectIOExplorer is a file explorer for supporting direct io.
type DirectIOExplorer struct {
	rootDir string
}

// NewDirectIOExplorer is a constructor of DirectIOExplorer.
func NewDirectIOExplorer(rootDir string) *DirectIOExplorer {
	return &DirectIOExplorer{
		rootDir: rootDir,
	}
}

// OpenFile opens file as direct io mode.
func (exp *DirectIOExplorer) OpenFile(filename domain.FileName) (*domain.File, error) {
	path := filepath.Join(exp.rootDir, string(filename))
	flag := os.O_RDWR | os.O_CREATE
	f, err := directio.OpenFile(path, flag, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return domain.NewFile(f), nil
}
