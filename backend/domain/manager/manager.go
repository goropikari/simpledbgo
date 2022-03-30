package manager

// type FileManager interface {
// 	OpenFile(filename FileName) (*File, error)
// 	CopyBlockToPage(block *Block, page *Page) error
// 	CopyPageToBlock(page *Page, block *Block) error
// 	AppendBlock(filename FileName) (*Block, error)
// 	NumBlock(filename FileName) (int, error)
// 	PrepagePage() (*Page, error)
// }
//
// // read(BlockId blk, Page p)
// // write(BlockId blk, Page p)
// // BlockId append(String filename)
// // int length(String filename)
// // int blockSize()
// // RandomAccessFile getFile(String filename)
//
// type LogManager interface {
// 	FlushCurrentPage()
// 	FlushByLSN(lsn int32)
// 	AppendRecord(record []byte)
// 	Iterator()
// }
//
// type BufferManager interface {
// 	Available() int
// 	FlushAll(txnum int32) error
// 	Pin(block *Block) (*Buffer, error)
// 	Unpin(*Buffer) error
// }
