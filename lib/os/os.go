package os

import (
	"os"

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

func OpenFile(path string, isDirectIO bool) (*File, error) {
	flag := os.O_RDWR | os.O_CREATE

	var f *os.File
	var err error

	if isDirectIO {
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

	return NewFile(f), nil
}
