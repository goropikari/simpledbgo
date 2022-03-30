package os

import (
	"os"
	"path/filepath"

	"github.com/goropikari/simpledb_go/backend/domain"
)

// NormalExplorer is a file explorer on normal mode.
type NormalExplorer struct {
	rootDir string
}

// NewNormalExplorer is a constructor of NewNormalExplorer.
func NewNormalExplorer(rootDir string) *NormalExplorer {
	return &NormalExplorer{
		rootDir: rootDir,
	}
}

// OpenFile opens a file.
func (exp *NormalExplorer) OpenFile(filename domain.FileName) (*domain.File, error) {
	path := filepath.Join(exp.rootDir, string(filename))
	flag := os.O_RDWR | os.O_CREATE
	f, err := os.OpenFile(path, flag, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return domain.NewFile(f), nil
}
