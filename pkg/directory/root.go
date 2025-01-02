package directory

import "github.com/vandi37/vanfs/pkg/files"

func NewRoot(path string) *Directory {
	root := &Directory{
		dirs:      map[string]*Directory{},
		files:     map[string]*files.File{},
		file_path: path,
	}
	root.last = root
	root.root = root
	return root
}
