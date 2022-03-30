package domain

import "os"

func (f *File) Remove() {
	os.Remove(string(f.Name()))
}
