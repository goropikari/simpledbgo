package bytes

import "github.com/goropikari/simpledb_go/lib/directio"

// DirectSliceCreater is a slice creature with direct io.
type DirectSliceCreater struct{}

// NewDirectSliceCreater is a constructor of DirectSliceCreater.
func NewDirectSliceCreater() *DirectSliceCreater {
	return &DirectSliceCreater{}
}

// Create creates a byte slice.
func (s *DirectSliceCreater) Create(n int) ([]byte, error) {
	b, err := directio.AlignedBlock(n)
	if err != nil {
		return nil, err
	}

	return b, nil
}
