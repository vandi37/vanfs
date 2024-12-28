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

func New(name string) (*File, error) {
	uniqueName := uuid.New().String() + "_" + time.Now().Format("2006.01.02_15.04.05") + "_" + name + ".vfs"

	file, err := os.Create(uniqueName)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorCreatingFile, err, vanerrors.EmptyHandler)
	}
	defer file.Close()

	return &File{
		path: uniqueName,
	}, nil
}

func (f *File) Remove() error {
	err := os.Remove(f.path)
	if err != nil {
		return vanerrors.NewWrap(ErrorRemovingFile, err, vanerrors.EmptyHandler)
	}
	return nil
}

func (f *File) Open() (*os.File, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorOpeningFile, err, vanerrors.EmptyHandler)
	}
	return file, nil
}

func (f *File) GetPath() string {
	return f.path
}
