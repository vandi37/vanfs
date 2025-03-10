package console

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vandi37/vanerrors"
	"github.com/vandi37/vanfs/pkg/cleaner"
	"github.com/vandi37/vanfs/pkg/filesystem"
)

type Console struct {
	fs      *filesystem.Filesystem
	cleaner cleaner.Cleaner
	funcs   map[string]func(s string) error
}

func New(fs *filesystem.Filesystem, cleaner cleaner.Cleaner) *Console {
	res := &Console{
		fs:      fs,
		cleaner: cleaner,
	}
	res.Init()
	return res
}

func (c *Console) Run(cancel context.CancelFunc) {
	for {
		var command string
		var data string
		path := c.fs.GetPath()
		if strings.HasPrefix(path, "/home") {
			path = "~" + path[5:]
		}
		fmt.Printf("\033[38;2;76;121;72m%s\033[0m:\033[38;2;255;255;255m%s\033[0m$ ", c.fs.Name, path)
		fmt.Scanln(&command, &data)
		var err error

		f, ok := c.funcs[command]
		switch {
		case ok:
			err = f(data)
		case command == "exit":
			cancel()
			return
		default:
			if command != "" {
				err = vanerrors.Simple("command not exist")
			}
		}
		if err != nil {
			fmt.Fprint(os.Stderr, "\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
		}
		err = c.fs.Save()
		if err != nil {
			fmt.Fprint(os.Stderr, "\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
		}
		err = c.fs.Reload()
		if err != nil {
			fmt.Fprint(os.Stderr, "\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
		}
		time.Sleep(time.Millisecond * 5)
	}
}
