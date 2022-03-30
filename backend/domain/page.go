package domain

// Page is a model of page.
type Page struct {
	ByteBuffer
}

// NewPage is a constructor of Page.
func NewPage(bb ByteBuffer) *Page {
	return &Page{
		ByteBuffer: bb,
	}
}
