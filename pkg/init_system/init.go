package init_system

import (
	"fmt"
	"os"
	"strings"
	"time"
	"vfs/pkg/filesystem"

	"github.com/AlecAivazis/survey/v2"
	"github.com/vandi37/vanerrors"
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

	go func() {
		time.Sleep(time.Millisecond)
		fmt.Printf("\n\033[0m\033[38;2;200;24;0m!!! Created system %s with backup file %s\n", systemName, backupPath)
	}()
	return fs, nil
}
