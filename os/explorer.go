package os

import (
	"log"
	"os"
	"path/filepath"

	"github.com/goropikari/simpledbgo/backend/domain"
)

type opener interface {
	OpenFile(string, int, os.FileMode) (*os.File, error)
}

// Explorer is a file explorer.
type Explorer struct {
	rootDir   string
	openFiles map[domain.FileName]*domain.File
	opener    opener
}

// NewExplorer is a constructor of NewExplorer.
func NewExplorer(rootDir string, opener opener) *Explorer {
	err := os.MkdirAll(rootDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return &Explorer{
		rootDir:   rootDir,
		openFiles: make(map[domain.FileName]*domain.File),
		opener:    opener,
	}
}

// OpenFile opens a file.
func (exp *Explorer) OpenFile(filename domain.FileName) (*domain.File, error) {
	if f, ok := exp.openFiles[filename]; ok {
		return f, nil
	}

	path := filepath.Join(exp.rootDir, string(filename))
	flag := os.O_RDWR | os.O_CREATE
	f, err := exp.opener.OpenFile(path, flag, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file := domain.NewFile(f)
	exp.openFiles[filename] = file

	return file, nil
}
