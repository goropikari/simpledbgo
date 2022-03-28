package os

import (
	"os"
	"path/filepath"

	"github.com/goropikari/simpledb_go/lib/directio"
)

type File struct {
	f *os.File
}

func NewFile(f *os.File) *File {
	return &File{
		f: f,
	}
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.f.Read(p)
}

func (f *File) Write(p []byte) (n int, err error) {
	return f.f.Write(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.f.Seek(offset, whence)
}

func (f *File) Close() error {
	return f.f.Close()
}

func (f *File) Size() (int64, error) {
	info, err := f.f.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

type NormalExplorer struct {
}

func NewNormalExplorer() *NormalExplorer {
	return &NormalExplorer{}
}

type DirectIOExplorer struct {
	NormalExplorer
}

func NewDirectIOExplorer() *DirectIOExplorer {
	return &DirectIOExplorer{}
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
