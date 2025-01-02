package filesystem

import "os"

func Init(name string, path string) (*Filesystem, error) {
	file, err := os.OpenFile(path+"tree.json", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	var fs = &Filesystem{
		Path:   path,
		Source: file,
		Name:   name,
	}
	fs.root = fs.Json.ToDir(path)
	fs.curDir = fs.root
	fs.Mkdir("/home")
	fs.Cd("/home")
	err = fs.Save()
	if err != nil {
		return nil, err
	}
	return fs, nil
}
