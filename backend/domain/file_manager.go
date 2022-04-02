package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

type FileManager interface {
	CopyBlockToPage(*Block, *Page) error
	CopyPageToBlock(*Page, *Block) error
	BlockLength(FileName) (int32, error)
	ExtendFile(FileName) (*Block, error)
	BlockSize() BlockSize
}
