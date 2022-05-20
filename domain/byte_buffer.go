package domain

import "io"

// ByteBuffer is an interface of buffer operations.
type ByteBuffer interface {
	io.ReadWriteSeeker
	GetData() []byte
	GetInt32(offset int64) (int32, error)
	GetUint32(offset int64) (uint32, error)
	GetString(offset int64) (string, error)
	GetBytes(offset int64) ([]byte, error)
	SetInt32(offset int64, val int32) error
	SetUint32(offset int64, val uint32) error
	SetString(offset int64, val string) error
	SetBytes(offset int64, val []byte) error
	Reset()
}
