package domain

import "io"

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// ByteBuffer is an interface of buffer operations.
type ByteBuffer interface {
	io.ReadWriteSeeker
	GetData() []byte
	GetInt32(offset int64) (int32, error)
	SetInt32(offset int64, val int32) error
	GetString(offset int64) (string, error)
	SetString(offset int64, val string) error
	GetBytes(offset int64) ([]byte, error)
	SetBytes(offset int64, val []byte) error
	NeededByteLength(x any) int64
	Size() int64
	Reset()
}
