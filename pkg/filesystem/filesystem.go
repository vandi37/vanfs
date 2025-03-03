package filesystem

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/vandi37/vanerrors"
	"github.com/vandi37/vanfs/pkg/directory"
)

const (
	ErrorToReadFile = "error to read file"
)

type Filesystem struct {
	root   *directory.Directory
	curDir *directory.Directory
	Path   string `json:"-"`
	Source interface {
		io.WriterAt
		io.ReadWriteCloser
		io.Seeker
	} `json:"-"`
	Name string            `json:"name"`
	Json directory.JsonDir `json:"files"`
}

func New(path string) (*Filesystem, error) {
	file, err := os.OpenFile(path+"tree.json", os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	var fs = &Filesystem{}

	err = json.NewDecoder(file).Decode(fs)

	if err != nil {
		return nil, err
	}

	fs.Source = file
	fs.root = fs.Json.ToDir(path)
	curDir, err := fs.root.OpenDirOrCreate("/home")
	if err != nil {
		return nil, err
	}
	fs.curDir = curDir
	fs.Path = path
	return fs, nil
}

func (f *Filesystem) String() string {
	return f.root.String()
}

func (f *Filesystem) Cd(path string) error {
	newDir, err := f.curDir.OpenDir(path)
	if err != nil {
		return err
	}
	f.curDir = newDir
	return nil
}

func (f *Filesystem) Of(path string) (*os.File, error) {
	return f.curDir.OpenFileOrAdd(path)
}

func (f *Filesystem) Cat(path string) (string, error) {
	file, err := f.curDir.OpenFile(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	res, err := io.ReadAll(file)
	if err != nil {
		return "", vanerrors.Simple(ErrorToReadFile)
	}
	return string(res), nil
}

func (f *Filesystem) Mkdir(path string) error {
	return f.curDir.AddDir(path)
}

func (f *Filesystem) Makefile(path string) error {
	return f.curDir.AddFile(path)
}

func (f *Filesystem) Rmdir(path string) error {
	return f.curDir.RemoveDir(path)
}

func (f *Filesystem) Rm(path string) error {
	return f.curDir.RemoveFile(path)
}

func (f *Filesystem) Save() error {
	f.Json = f.root.ToJsonDir()
	res, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Source.WriteAt(res, 0)
	return err
}

func (f *Filesystem) GetPath() string {
	return f.curDir.GetPath()
}

func (f *Filesystem) Ls(path string) ([]string, error) {
	lsDir := f.curDir
	if path != "" {
		var err error
		lsDir, err = f.curDir.OpenDir(path)
		if err != nil {
			return nil, err
		}
	}

	return lsDir.List(), nil
}

func (f *Filesystem) Tree(path string) (string, error) {
	treeDir := f.curDir
	if path != "" {
		var err error
		treeDir, err = f.curDir.OpenDir(path)
		if err != nil {
			return "", err
		}
	}

	return treeDir.String(), nil
}

func (f *Filesystem) Reload() error {
	path := f.GetPath()
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	d, err := f.root.OpenDir(path)
	if err != nil {
		return err
	}
	f.curDir = d
	return nil
}
