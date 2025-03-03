package files

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/vandi37/vanerrors"
)

const (
	ErrorCreatingFile = "error creating file"
	ErrorRemovingFile = "error removing file"
	ErrorOpeningFile  = "error opening file"
)

type File struct {
	path string
}

func Create(name string, path string) (*File, error) {
	uniqueName := path + uuid.New().String() + "_" + time.Now().Format("2006.01.02_15.04.05") + "_" + name + ".vfs"

	file, err := os.OpenFile(uniqueName, os.O_CREATE, 0777)
	if err != nil {
		return nil, vanerrors.Wrap(ErrorCreatingFile, err)
	}
	defer file.Close()

	return &File{
		path: uniqueName,
	}, nil
}

func (f *File) Remove() error {
	err := os.Remove(f.path)
	if err != nil {
		return vanerrors.Wrap(ErrorRemovingFile, err)
	}
	return nil
}

func (f *File) Open() (*os.File, error) {
	file, err := os.OpenFile(f.path, os.O_RDWR, 0777)
	if err != nil {
		return nil, vanerrors.Wrap(ErrorOpeningFile, err)
	}
	return file, nil
}

func (f *File) GetPath() string {
	return f.path
}

func New(path string) *File {
	return &File{path: path}
}
