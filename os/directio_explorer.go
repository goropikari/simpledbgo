package os

import (
	"log"
	"os"
	"path/filepath"

	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/lib/directio"
)

// DirectIOExplorer is a file explorer for supporting direct io.
type DirectIOExplorer struct {
	rootDir   string
	openFiles map[domain.FileName]*domain.File
}

// NewDirectIOExplorer is a constructor of DirectIOExplorer.
func NewDirectIOExplorer(rootDir string) *DirectIOExplorer {
	err := os.MkdirAll(rootDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return &DirectIOExplorer{
		rootDir:   rootDir,
		openFiles: make(map[domain.FileName]*domain.File),
	}
}

// OpenFile opens file as direct io mode.
func (exp *DirectIOExplorer) OpenFile(filename domain.FileName) (*domain.File, error) {
	if f, ok := exp.openFiles[filename]; ok {
		return f, nil
	}

	path := filepath.Join(exp.rootDir, string(filename))
	flag := os.O_RDWR | os.O_CREATE
	f, err := directio.OpenFile(path, flag, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file := domain.NewFile(f)
	exp.openFiles[filename] = file

	return file, nil
}
