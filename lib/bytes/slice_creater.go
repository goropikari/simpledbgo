package bytes

// SliceCreater is a model of creature of byte slice.
type SliceCreater struct{}

// NewSliceCreater is a constructor of SliceCreater.
func NewSliceCreater() *SliceCreater {
	return &SliceCreater{}
}

// Create creats an n bytes `bytes` slice.
func (s *SliceCreater) Create(n int) ([]byte, error) {
	return make([]byte, n), nil
}
