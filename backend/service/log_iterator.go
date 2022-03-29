package service

type Iterator interface {
	HasNext() bool
	Next() ([]byte, error)
}
