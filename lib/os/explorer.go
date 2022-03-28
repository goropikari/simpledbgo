package os

import (
	"os"
	"path/filepath"

	"github.com/goropikari/simpledb_go/lib/directio"
)

type NormalExplorer struct{}

func NewNormalExplorer() *NormalExplorer {
	return &NormalExplorer{}
}

type DirectIOExplorer struct {
	*NormalExplorer
}

func NewDirectIOExplorer() *DirectIOExplorer {
	return &DirectIOExplorer{
		NormalExplorer: NewNormalExplorer(),
	}
}

func (ex *NormalExplorer) MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func (ex *NormalExplorer) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (ex *NormalExplorer) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (ex *NormalExplorer) Remove(dir string, file string) error {
	return os.Remove(filepath.Join(dir, file))
}

func (ex *NormalExplorer) OpenFile(path string) (*File, error) {
	flag := os.O_RDWR | os.O_CREATE

	f, err := os.OpenFile(path, flag, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return NewFile(f), nil
}

func (ex *DirectIOExplorer) OpenFile(path string) (*File, error) {
	flag := os.O_RDWR | os.O_CREATE

	f, err := directio.OpenFile(path, flag, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return NewFile(f), nil
}
