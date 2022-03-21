package service

import "github.com/goropikari/simpledb_go/backend/core"

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/tests/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// FileManager is an interface of file manager.
type FileManager interface {
	// GetBlockSize returns block size.
	GetBlockSize() int

	// CopyBlockToPage copies block to page.
	CopyBlockToPage(block *core.Block, page *core.Page) error

	// CopyPageToBlock copies page to block.
	CopyPageToBlock(page *core.Page, block *core.Block) error

	// AppendBlock appends block file.
	AppendBlock(filename core.FileName) (*core.Block, error)

	// LastBlock returns last block of the file.
	LastBlock(filename core.FileName) (*core.Block, error)

	// PreparePage prepares a page.
	PreparePage() (*core.Page, error)

	// IsZero checks whether given file manager is zero value or not.
	IsZero() bool
}
