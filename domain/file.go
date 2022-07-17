package domain

import (
	"io"
	"os"

	"github.com/goropikari/simpledbgo/errors"
	"github.com/goropikari/simpledbgo/lib/bytes"
)

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// Explorer is an interface of file explorer.
type Explorer interface {
	OpenFile(FileName) (*File, error)
}

// ByteSliceFactory is a factory of byte slice.
type ByteSliceFactory interface {
	Create(int) ([]byte, error)
}

/*
	File
*/

// ErrInvalidFileName is an error that means invalid file name.
var ErrInvalidFileName = errors.New("invalid file name")

// FileName is a value object of file name.
type FileName string

// NewFileName is a constructor of FileName.
func NewFileName(name string) (FileName, error) {
	if name == "" {
		return "", ErrInvalidFileName
	}

	return FileName(name), nil
}

// String stringfies file name.
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
func (f *File) Seek(offset int64) (int64, error) {
	return f.file.Seek(offset, io.SeekStart)
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

/*
	Block
*/
var (
	// ErrNegativeBlockNumber means given block number is non negative.
	ErrNegativeBlockNumber = errors.New("block number must be non negative")

	// ErrNonPositiveBlockSize means given block size must be positive.
	ErrNonPositiveBlockSize = errors.New("block size must be positive")
)

// BlockNumber is value object of block number.
type BlockNumber int32

// NewBlockNumber is a constructor of BlockNumber.
func NewBlockNumber(n int32) (BlockNumber, error) {
	if n < 0 {
		return 0, ErrNegativeBlockNumber
	}

	return BlockNumber(n), nil
}

func (bn BlockNumber) ToInt32() int32 {
	return int32(bn)
}

// BlockSize is value object of block size.
type BlockSize int32

// NewBlockSize is a constructor of BlockSize.
func NewBlockSize(n int32) (BlockSize, error) {
	if n <= 0 {
		return 0, ErrNonPositiveBlockSize
	}

	return BlockSize(n), nil
}

// Block is a model of block.
type Block struct {
	filename FileName
	number   BlockNumber
}

// NewBlock is a constructor of Block.
func NewBlock(filename FileName, number BlockNumber) Block {
	return Block{
		filename: filename,
		number:   number,
	}
}

// NewDummyBlock constructs a dummy Block.
func NewDummyBlock(filename FileName) Block {
	return Block{
		filename: filename,
		number:   -1,
	}
}

// Equal compares equality of two blocks.
func (b Block) Equal(other Block) bool {
	return b == other
}

// FileName returns corresponding file name.
func (b Block) FileName() FileName {
	return b.filename
}

// Number returns block number.
func (b Block) Number() BlockNumber {
	return b.number
}

/*
	Page
*/

// Page is a model of page.
type Page struct {
	ByteBuffer
}

// NewPage is a constructor of Page.
func NewPage(bb ByteBuffer) *Page {
	return &Page{
		ByteBuffer: bb,
	}
}

// PageFactory is a factory of page.
type PageFactory struct {
	bsf       ByteSliceFactory
	blockSize BlockSize
}

// NewPageFactory is a constructor of PageFactory.
func NewPageFactory(bsf ByteSliceFactory, blockSize BlockSize) *PageFactory {
	return &PageFactory{
		bsf:       bsf,
		blockSize: blockSize,
	}
}

// Create creates a page.
func (pf *PageFactory) Create() (*Page, error) {
	b, err := pf.bsf.Create(int(pf.blockSize))
	if err != nil {
		return nil, errors.Err(err, "create byte slice")
	}

	bb := bytes.NewBufferBytes(b)
	page := NewPage(bb)

	return page, nil
}
