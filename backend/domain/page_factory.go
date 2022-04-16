package domain

import "github.com/goropikari/simpledbgo/lib/bytes"

// PageFactory is a factory of page.
type PageFactory struct {
	bsf       ByteSliceFactory
	blockSize BlockSize
}

// NewPageFactory is a constructor of PageFactory.
func NewPageFactory(bsf ByteSliceFactory, blockSize BlockSize) *PageFactory {
	return &PageFactory{
		bsf:       bsf,
		blockSize: blockSize,
	}
}

// Create creates a page.
func (pf *PageFactory) Create() (*Page, error) {
	b, err := pf.bsf.Create(int(pf.blockSize))
	if err != nil {
		return nil, err
	}

	bb := bytes.NewBufferBytes(b)
	page := NewPage(bb)

	return page, nil
}
