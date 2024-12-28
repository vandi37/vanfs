package directory

import (
	"strings"
	"vfs/pkg/files"

	"github.com/vandi37/vanerrors"
)

const (
	EmptyPath              = "empty path"
	DirectoryExists        = "directory exists"
	DirectoryDoesNotExists = "directory does not exists"
	FileExists             = "file exists"
)

type Directory struct {
	dirs  map[string]*Directory
	files map[string]*files.File
	last  *Directory
	root  *Directory
	name  string
}

func (d *Directory) OpenDir(path string) (*Directory, error) {
	return d.openDir(path, true)
}

func (d *Directory) openDir(path string, autoCreate bool) (*Directory, error) {
	if len(path) < 1 {
		return nil, vanerrors.NewSimple(EmptyPath)
	}

	if path[0] == '/' {
		return d.root.openDir(path[:1], autoCreate)
	}

	paths := strings.Split(path, "/")

	currentDir := d
	for _, p := range paths {
		if p == ".." {
			currentDir = currentDir.last
			continue
		}

		dir, ok := currentDir.dirs[p]
		if !autoCreate && (!ok || dir == nil) {
			return nil, vanerrors.NewSimple(DirectoryDoesNotExists, p)
		} else if autoCreate && (!ok || dir == nil) {
			currentDir.dirs[p] = &Directory{
				dirs:  map[string]*Directory{},
				files: map[string]*files.File{},
				last:  currentDir,
				root:  currentDir.root,
				name:  p,
			}
			dir = currentDir.dirs[p]
		}
		currentDir = dir
	}
	return currentDir, nil
}

func (d *Directory) addDir(path string, errorIfExist bool) error {
	if len(path) < 1 {
		return vanerrors.NewSimple(EmptyPath)
	}
	if path[0] == '/' {
		return d.root.addDir(path[:1], errorIfExist)
	}
	paths := strings.Split(path, "/")

	currentDir := d
	if len(paths) > 1 {
		var err error
		currentDir, err = d.openDir(strings.Join(paths[1:], "/"), true)
		if err != nil {
			return err
		}
	}

	dir, ok := currentDir.dirs[paths[len(paths)-1]]
	if errorIfExist && ok && dir != nil {
		return vanerrors.NewSimple(DirectoryExists, paths[len(paths)-1])
	} else if !ok || dir == nil {
		currentDir.dirs[paths[len(paths)-1]] = &Directory{
			dirs:  map[string]*Directory{},
			files: map[string]*files.File{},
			last:  d,
			root:  d.root,
		}
	}
	return nil
}

func (d *Directory) AddDir(path string) error {
	return d.addDir(path, true)
}

func (d *Directory) AddDirIfNotExists(path string) error {
	return d.addDir(path, false)
}

func (d *Directory) addFile(path string, errorIfExist bool) error {
	if len(path) < 1 {
		return vanerrors.NewSimple(EmptyPath)
	}
	if path[0] == '/' {
		return d.root.addFile(path[:1], errorIfExist)
	}

	paths := strings.Split(path, "/")

	var last = len(paths) - 1

	currentDir := d
	if len(paths) > 1 {
		var err error
		currentDir, err = d.openDir(strings.Join(paths[:last], "/"), true)
		if err != nil {
			return err
		}
	}

	f, ok := currentDir.files[paths[last]]
	if errorIfExist && ok && f != nil {
		return vanerrors.NewSimple(FileExists, paths[last])
	} else if !ok || f == nil {
		file, err := files.New(paths[last])
		if err != nil {
			return err
		}
		currentDir.files[paths[last]] = file
	}
	return nil
}

func (d *Directory) AddFile(path string) error {
	return d.addFile(path, true)
}

func (d *Directory) AddFileIfNotExists(path string) error {
	return d.addFile(path, false)
}

func (d *Directory) IsEmpty() bool {
	for _, d := range d.dirs {
		if d != nil {
			return false
		}
	}

	for _, f := range d.files {
		if f != nil {
			return false
		}

	}
	return true
}

func (d *Directory) selfRemove() error {
	for n, f := range d.files {
		err := f.Remove()
		if err != nil {
			return err
		}
		delete(d.files, n)
	}
	for n, dir := range d.dirs {
		err := dir.selfRemove()
		if err != nil {
			return err
		}
		delete(d.dirs, n)
	}
	return nil
}

func (d *Directory) removeDir(path string, errorIfNotExist bool) error {
	if len(path) < 1 {
		if d == d.root {
			return d.selfRemove()
		}
		d.last.removeDir(d.name, errorIfNotExist)
	}
	if path[0] == '/' {
		return d.root.removeDir(path[:1], errorIfNotExist)
	}
	paths := strings.Split(path, "/")

	currentDir := d
	if len(paths) > 1 {
		var err error
		currentDir, err = d.openDir(strings.Join(paths[1:], "/"), false)
		if err != nil {
			return err
		}
	}

	dir, ok := currentDir.dirs[paths[len(paths)-1]]
	if errorIfNotExist && ok && dir != nil {
		return vanerrors.NewSimple(DirectoryDoesNotExists, paths[len(paths)-1])
	} else if !ok || dir == nil {
		currentDir.selfRemove()
		delete(currentDir.dirs, paths[len(paths)-1])
	}
	return nil
}

func (d *Directory) RemoveDir(path string) error {
	return d.removeDir(path, true)
}

func (d *Directory) RemoveDirIfExists(path string) error {
	return d.removeDir(path, false)
}
