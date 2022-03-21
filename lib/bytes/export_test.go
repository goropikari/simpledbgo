package bytes

// GetBuf returns buffer for test.
func (buf *Buffer) GetBuf() []byte {
	return buf.buf
}

// GetOff returns offset for test.
func (buf *Buffer) GetOff() int64 {
	return buf.off
}
