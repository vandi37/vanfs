package init_system

import (
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/vandi37/vanerrors"
	"github.com/vandi37/vanfs/pkg/filesystem"
)

const (
	ErrorSettingEnv  = "error setting env"
	ErrorGettingPath = "error getting path"
)

func Init() (*filesystem.Filesystem, error) {
	var systemName string = "vfs"
	var backupPath string

	prompt1 := &survey.Input{
		Message: "System name:",
	}
	survey.AskOne(prompt1, &systemName)

	prompt2 := &survey.Input{
		Message: "Path to backup file:",
	}
	survey.AskOne(prompt2, &backupPath)

	if systemName == "" {
		systemName = "vfs"
	}
	if backupPath == "" {
		path, err := os.Getwd()
		if err != nil {
			return nil, vanerrors.NewWrap(ErrorGettingPath, err, vanerrors.EmptyHandler)
		}
		backupPath = path + "/vfs_backup/"
		os.Mkdir(backupPath, 0777)
	}

	if !strings.HasSuffix(backupPath, "/") {
		backupPath += "/"
	}

	fs, err := filesystem.Init(systemName, backupPath)
	if err != nil {
		return nil, err
	}
	return fs, nil
}
