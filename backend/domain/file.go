package domain

import (
	"os"
)

// FileName is a value object of file name.
type FileName string

// NewFileName is a constructor of FileName.
func NewFileName(name string) (FileName, error) {
	if name == "" {
		return "", ErrInvalidFileName
	}

	return FileName(name), nil
}

// String stringfy file name.
func (f FileName) String() string {
	return string(f)
}

// File is a model of file.
type File struct {
	file *os.File
}

// NewFile is a constructor of File.
func NewFile(f *os.File) *File {
	return &File{
		file: f,
	}
}

// Read reads up to len(b) bytes from the File and stores them in b.
func (f *File) Read(b []byte) (n int, err error) {
	return f.file.Read(b)
}

// Write writes len(b) bytes from b to the File.
func (f *File) Write(b []byte) (n int, err error) {
	return f.file.Write(b)
}

// Seek sets the offset for the next Read or Write on file to offset.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

// Close closes the File.
func (f *File) Close() error {
	return f.file.Close()
}

// Size returns the size of file.
func (f *File) Size() (int64, error) {
	info, err := f.file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

// Name returns the file name.
func (f *File) Name() FileName {
	name := f.file.Name()
	filename, _ := NewFileName(name)

	return filename
}
