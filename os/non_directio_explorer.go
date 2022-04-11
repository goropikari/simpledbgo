package os

import (
	"os"
)

// NonDirectIOExplorer is a file explorer on normal mode.
type NonDirectIOExplorer struct {
	*Explorer
}

// NewNonDirectIOExplorer is a constructor of NewNonDirectIOExplorer.
func NewNonDirectIOExplorer(rootDir string) *NonDirectIOExplorer {
	opener := newNonDirectIOOpener()
	explorer := NewExplorer(rootDir, opener)

	return &NonDirectIOExplorer{
		Explorer: explorer,
	}
}

// nonDirectIOOpener is opener for non direct io.
type nonDirectIOOpener struct{}

func newNonDirectIOOpener() *nonDirectIOOpener {
	return &nonDirectIOOpener{}
}

// OpenFile opens a file.
func (op *nonDirectIOOpener) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}
