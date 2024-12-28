package directory

import "vfs/pkg/files"

func NewRoot() *Directory {
	root := &Directory{
		dirs:  map[string]*Directory{},
		files: map[string]*files.File{},
	}
	root.last = root
	root.root = root
	return root
}
