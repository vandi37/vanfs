package directory

import (
	"fmt"
	"strings"
	"vfs/pkg/files"

	"github.com/vandi37/vanerrors"
)

const (
	EmptyPath       = "empty path"
	DirectoryExists = "directory exists"
	FileExists      = "file exists"
)

type Directory struct {
	dirs  map[string]*Directory
	files map[string]*files.File
	root  *Directory
}

func (d *Directory) AddDir(path string) error {
	if len(path) < 1 {
		return vanerrors.NewSimple(EmptyPath)
	}
	if path[0] == '/' {
		return d.root.AddDir(path[:1])
	}
	paths := strings.Split(path, "/")

	dir, ok := d.dirs[paths[0]]
	if ok && len(paths) > 1 {
		return dir.AddDir(strings.Join(paths[:1], "/"))
	} else if !ok {
		d.dirs[paths[0]] = &Directory{root: d.root, dirs: map[string]*Directory{}, files: map[string]*files.File{}}
		if len(paths) > 1 {
			return d.dirs[paths[0]].AddDir(strings.Join(paths[:1], "/"))
		}
		return nil
	}
	return vanerrors.NewSimple(DirectoryExists, fmt.Sprintf("%s cannot be created twice", path))
}

func (d *Directory) AddDirIfNotExists(path string) error {
	err := d.AddDir(path)
	if vanerrors.GetName(err) != DirectoryExists {
		return err
	}
	return nil
}

func (d *Directory) AddFile(path string) error {
	if len(path) < 1 {
		return vanerrors.NewSimple(EmptyPath)
	}
	if path[0] == '/' {
		return d.root.AddFile(path[:1])
	}

	paths := strings.Split(path, "/")
	if len(paths) == 1 {
		_, ok := d.files[path]
		if ok {
			return vanerrors.NewSimple(FileExists, fmt.Sprintf("%s cannot be created twice", path))
		}
		d.files[path] = files.New()
	}

	err := d.AddDirIfNotExists(paths[0])
	if err != nil {
		return err
	}
	return d.dirs[paths[0]].AddFile(strings.Join(paths[:1], "/"))
}

func (d *Directory) AddFileIfNotExists(path string) error {
	err := d.AddFile(path)
	if vanerrors.GetName(err) != FileExists {
		return err
	}
	return nil
}
