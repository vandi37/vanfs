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
	FileDoesNotExists      = "file does not exists"
)

type Directory struct {
	dirs  map[string]*Directory
	files map[string]*files.File
	last  *Directory
	root  *Directory
	name  string
}

func (d *Directory) OpenDir(path string) (*Directory, error) {
	return d.openDir(path, false)
}

func (d *Directory) OpenDirOrCreate(path string) (*Directory, error) {
	return d.openDir(path, true)
}

func (d *Directory) openDir(path string, autoCreate bool) (*Directory, error) {
	if len(path) < 1 {
		return nil, vanerrors.NewSimple(EmptyPath)
	}

	if path[len(path)-1] == '/' && path != "/" {
		path = path[1:]
	}

	if path == "/" {
		return d.root, nil
	}

	if path[0] == '/' {
		return d.root.openDir(path[1:], autoCreate)
	}
	paths := strings.Split(path, "/")

	currentDir := d
	for _, p := range paths {
		if p == ".." {
			if currentDir != currentDir.last {
				currentDir = currentDir.last
			}
			continue
		}

		dir, ok := currentDir.dirs[p]
		if !autoCreate && (!ok || dir == nil) {
			return nil, vanerrors.NewSimple(DirectoryDoesNotExists, p)
		} else if !ok || dir == nil {
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

	if path[len(path)-1] == '/' && path != "/" {
		path = path[1:]
	}

	if path[0] == '/' {
		return d.root.addDir(path[1:], errorIfExist)
	}
	paths := strings.Split(path, "/")

	currentDir := d
	if len(paths) > 1 {
		var err error
		currentDir, err = d.OpenDirOrCreate(strings.Join(paths[:len(paths)-1], "/"))
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
			name:  paths[len(paths)-1],
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
		return d.selfRemove()
	}

	if path[len(path)-1] == '/' && path != "/" {
		path = path[1:]
	}

	if path[0] == '/' {
		return d.root.removeDir(path[1:], errorIfNotExist)
	}

	paths := strings.Split(path, "/")

	currentDir := d
	if len(paths) > 1 {
		var err error
		currentDir, err = d.OpenDir(strings.Join(paths[:len(paths)-1], "/"))
		if vanerrors.GetName(err) == DirectoryDoesNotExists && !errorIfNotExist {
			return nil
		}
		if err != nil {
			return err
		}
	}

	dir, ok := currentDir.dirs[paths[len(paths)-1]]
	if errorIfNotExist && (!ok || dir == nil) {
		return vanerrors.NewSimple(DirectoryDoesNotExists, paths[len(paths)-1])
	} else if ok {
		err :=dir.selfRemove()
		if err != nil {
			return err
		}
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

func (d *Directory) GetPath() string {
	var res string
	currentDir := d
	for currentDir != d.root {
		res = currentDir.name + "/" + res
		currentDir = currentDir.last
	}
	res = "/" + res
	return res
}

func (d *Directory) List() []string {
	var res = []string{}

	for n := range d.files {
		res = append(res, n)
	}
	for n := range d.dirs {
		res = append(res, n)
	}
	return res
}
