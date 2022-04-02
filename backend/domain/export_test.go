package domain

func (buf *Buffer) LSN() LSN {
	return buf.lsn
}

func (buf *Buffer) Page() *Page {
	return buf.page
}
