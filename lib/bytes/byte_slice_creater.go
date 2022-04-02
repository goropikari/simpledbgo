package bytes

// ByteSliceCreater is a model of creature of byte slice.
type ByteSliceCreater struct{}

// NewByteSliceCreater is a constructor of ByteSliceCreater.
func NewByteSliceCreater() *ByteSliceCreater {
	return &ByteSliceCreater{}
}

// Create creats an n bytes `bytes` slice.
func (s *ByteSliceCreater) Create(n int) ([]byte, error) {
	return make([]byte, n), nil
}
