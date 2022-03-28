package os

import (
	"os"
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
