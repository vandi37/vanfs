package directory

import (
	"errors"
	"os"
	"strings"

	"github.com/vandi37/vanerrors"
	"github.com/vandi37/vanfs/pkg/files"
)

func (d *Directory) addFile(path string, errorIfExist bool) error {
	if len(path) < 1 {
		return vanerrors.Simple(EmptyPath)
	}

	if path[len(path)-1] == '/' && path != "/" {
		path = path[1:]
	}

	if path[0] == '/' {
		return d.root.addFile(path[1:], errorIfExist)
	}
	paths := strings.Split(path, "/")

	var last = len(paths) - 1

	currentDir := d
	if len(paths) > 1 {
		var err error
		currentDir, err = d.OpenDirOrCreate(strings.Join(paths[:last], "/"))
		if err != nil {
			return err
		}
	}

	f, ok := currentDir.files[paths[last]]
	if errorIfExist && ok && f != nil {
		return vanerrors.New(FileExists, paths[last])
	} else if !ok || f == nil {
		file, err := files.Create(paths[last], currentDir.file_path)
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

func (d *Directory) removeFile(path string, errorIfNotExist bool) error {
	if len(path) < 1 {
		return vanerrors.Simple(EmptyPath)
	}
	if path[len(path)-1] == '/' && path != "/" {
		path = path[1:]
	}

	if path[0] == '/' {
		return d.root.removeFile(path[1:], errorIfNotExist)
	}
	paths := strings.Split(path, "/")

	currentDir := d
	if len(paths) > 1 {
		var err error
		currentDir, err = d.OpenDir(strings.Join(paths[:len(paths)-1], "/"))
		if errors.Is(err, vanerrors.Simple(FileDoesNotExists)) && !errorIfNotExist {
			return nil
		}
		if err != nil {
			return err
		}
	}

	file, ok := currentDir.files[paths[len(paths)-1]]
	if errorIfNotExist && (!ok || file == nil) {
		return vanerrors.New(FileDoesNotExists, paths[len(paths)-1])
	} else if ok && file != nil {
		err := file.Remove()
		if err != nil {
			return err
		}
		delete(currentDir.files, paths[len(paths)-1])
	}
	return nil
}
func (d *Directory) RemoveFile(path string) error {
	return d.removeFile(path, true)
}

func (d *Directory) RemoveFileIfExists(path string) error {
	return d.removeFile(path, false)
}

func (d *Directory) openFile(path string, errorIfNotExist bool) (*os.File, error) {
	if len(path) < 1 {
		return nil, vanerrors.Simple(EmptyPath)
	}

	if path[len(path)-1] == '/' && path != "/" {
		path = path[1:]
	}

	if path[0] == '/' {
		return d.root.openFile(path[1:], errorIfNotExist)
	}
	paths := strings.Split(path, "/")

	currentDir := d
	if len(paths) > 1 {
		var err error
		currentDir, err = d.OpenDirOrCreate(strings.Join(paths[:len(paths)-1], "/"))
		if err != nil {
			return nil, err
		}
	}

	file, ok := currentDir.files[paths[len(paths)-1]]
	if errorIfNotExist && (!ok || file == nil) {
		return nil, vanerrors.New(FileDoesNotExists, paths[len(paths)-1])
	} else if !ok || file == nil {
		err := currentDir.AddFile(paths[len(paths)-1])
		if err != nil {
			return nil, err
		}
		file = currentDir.files[paths[len(paths)-1]]
	}
	return file.Open()
}

func (d *Directory) OpenFile(path string) (*os.File, error) {
	return d.openFile(path, true)
}

func (d *Directory) OpenFileOrAdd(path string) (*os.File, error) {
	return d.openFile(path, false)

}
