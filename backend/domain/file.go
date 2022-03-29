package domain

import (
	"os"
)

type FileName string

func NewFileName(name string) (FileName, error) {
	if name == "" {
		return "", ErrInvalidFileName
	}

	return FileName(name), nil
}

func (f FileName) String() string {
	return string(f)
}

type File struct {
	file *os.File
}

func NewFile(f *os.File) *File {
	return &File{
		file: f,
	}
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *File) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) Size() (int64, error) {
	info, err := f.file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func (f *File) Name() FileName {
	name := f.file.Name()
	filename, _ := NewFileName(name)

	return filename
}
