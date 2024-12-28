package directory

import "encoding/json"

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
