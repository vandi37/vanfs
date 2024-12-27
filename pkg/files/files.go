package files

import "bytes"

type File struct {
	*bytes.Buffer
}

func New() *File {
	return &File{bytes.NewBuffer([]byte{})}
}
