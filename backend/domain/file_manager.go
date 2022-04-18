package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// FileManager is an interface of file manager.
type FileManager interface {
	CopyBlockToPage(Block, *Page) error
	CopyPageToBlock(*Page, Block) error
	BlockLength(FileName) (int32, error)
	ExtendFile(FileName) (Block, error)
	BlockSize() BlockSize
}
