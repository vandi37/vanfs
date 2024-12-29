package directory

import (
	"encoding/json"
	"vfs/pkg/files"
)

type JsonDir struct {
	Dirs  map[string]JsonDir `json:"dirs"`
	Files map[string]string  `json:"files"`
}

func (d *Directory) ToJsonDir() JsonDir {
	res := JsonDir{
		Dirs:  map[string]JsonDir{},
		Files: map[string]string{},
	}

	for path, file := range d.files {
		res.Files[path] = file.GetPath()
	}
	for path, dir := range d.dirs {
		res.Dirs[path] = dir.ToJsonDir()
	}
	return res
}

func (d JsonDir) ToJson() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

func (d JsonDir) ToDir() *Directory {
	var res = NewRoot()
	res.addJsonDir(d)
	return res
}

func (d *Directory) addJsonDir(j JsonDir) {
	for path, file := range j.Files {
		d.files[path] = files.New(file)
	}
	for path, dir := range j.Dirs {
		d.dirs[path] = &Directory{
			dirs:  map[string]*Directory{},
			files: map[string]*files.File{},
			last:  d,
			root:  d.root,
			name:  path,
		}
		d.dirs[path].addJsonDir(dir)
	}
}
