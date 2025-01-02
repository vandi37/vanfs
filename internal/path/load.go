package path

import (
	"os"
	"strings"

	"github.com/vandi37/vanerrors"
)

const (
	PathNotFound = "path not found"
)

const (
	VFS_PATH = "VFS_PATH"
)

func LoadPath() (string, error) {
	path := os.Getenv(VFS_PATH)

	if path == "" {
		return "", vanerrors.NewSimple(PathNotFound)
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path, nil
}
