package load

import (
	"fmt"

	"github.com/vandi37/flags"
	"github.com/vandi37/vanerrors"
)

type Loading struct {
	Help     bool   `flags:"help"`
	New      bool   `flags:"new"`
	Path     string `flags:"path"`
	BoolPath bool   `flags:"path"`
	Default  bool   `falags:"default"`
}

var shortcuts = map[rune]string{
	'H': "help",
	'N': "new",
	'P': "path",
	'D': "default",
}

const (
	PathHasNoArguments = "path has no arguments"
)

func Load() (int, bool, error) {
	load := new(Loading)
	if err := flags.ArgsWithShortcuts(load, shortcuts); err != nil {
		return -1, false, err
	}

	// Help
	if load.Help {
		fmt.Println(load.ProseccHelp())
		return -1, true, nil
	}

	if load.BoolPath {
		return -1, false, vanerrors.Simple(PathHasNoArguments)
	}

	if load.New && !load.Default && load.Path == "" {
		return 2, false, nil
	}

	if load.Path != "" && !load.New && !load.Default {
		return 1, false, nil
	}

	if load.Default && !load.New && load.Path == "" {
		return 0, false, nil
	}

	return -1, false, nil
}

var help = map[string]string{
	"new":     `--new (-N) Creates a new vanfs. Step by step: System name, optional: Path to backup file`,
	"path":    `--path "path" (-P "path") Loading from path to backup file.`,
	"default": `--default (-D) Loading from env 'VFS_PATH. `,
}

func (l *Loading) ProseccHelp() string {
	needHelp := []string{}
	if l.New {
		needHelp = append(needHelp, help["new"])
	}
	if l.BoolPath {
		needHelp = append(needHelp, help["path"])
	}

	if l.Default {
		needHelp = append(needHelp, help["default"])
	}

	var res string

	for i, advice := range needHelp {
		res += advice
		if i != len(needHelp)-1 {
			res += "\n\n"
		}
	}
	return res
}
