package bytes

import "github.com/goropikari/simpledb_go/lib/directio"

type DirectSliceCreater struct{}

func NewDirectSliceCreater() *DirectSliceCreater {
	return &DirectSliceCreater{}
}

func (s *DirectSliceCreater) Create(n int) ([]byte, error) {
	b, err := directio.AlignedBlock(n)
	if err != nil {
		return nil, err
	}

	return b, nil
}
