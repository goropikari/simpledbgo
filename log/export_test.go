package log

// FlushByLSN exports flushByLSN for test.
func (mgr *Manager) FlushByLSN(x int) error {
	return mgr.flushByLSN(x)
}
