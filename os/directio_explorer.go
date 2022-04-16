package os

import (
	"os"

	"github.com/goropikari/simpledbgo/lib/directio"
)

// DirectIOExplorer is a file explorer for supporting direct io.
type DirectIOExplorer struct {
	*Explorer
}

// NewDirectIOExplorer is a constructor of DirectIOExplorer.
func NewDirectIOExplorer(rootDir string) *DirectIOExplorer {
	opener := newDirectIOOpener()
	explorer := NewExplorer(rootDir, opener)

	return &DirectIOExplorer{
		Explorer: explorer,
	}
}

type directIOOpener struct{}

func newDirectIOOpener() *directIOOpener {
	return &directIOOpener{}
}

// OpenFile opens a file.
func (op *directIOOpener) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return directio.OpenFile(name, flag, perm)
}
