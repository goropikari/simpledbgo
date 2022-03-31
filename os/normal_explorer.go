package os

import (
	"os"
	"path/filepath"

	"github.com/goropikari/simpledb_go/backend/domain"
)

// NormalExplorer is a file explorer on normal mode.
type NormalExplorer struct {
	rootDir   string
	openFiles map[domain.FileName]*domain.File
}

// NewNormalExplorer is a constructor of NewNormalExplorer.
func NewNormalExplorer(rootDir string) *NormalExplorer {
	return &NormalExplorer{
		rootDir:   rootDir,
		openFiles: make(map[domain.FileName]*domain.File),
	}
}

// OpenFile opens a file.
func (exp *NormalExplorer) OpenFile(filename domain.FileName) (*domain.File, error) {
	if f, ok := exp.openFiles[filename]; ok {
		return f, nil
	}

	path := filepath.Join(exp.rootDir, string(filename))
	flag := os.O_RDWR | os.O_CREATE
	f, err := os.OpenFile(path, flag, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file := domain.NewFile(f)
	exp.openFiles[filename] = file

	return file, nil
}
