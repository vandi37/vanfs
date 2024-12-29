package console

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"vfs/pkg/cleaner"
	"vfs/pkg/filesystem"
	"vfs/pkg/ide"

	"github.com/vandi37/vanerrors"
)

type Console struct {
	fs      *filesystem.Filesystem
	cleaner cleaner.Cleaner
}

func New(fs *filesystem.Filesystem, cleaner cleaner.Cleaner) *Console {
	return &Console{
		fs:      fs,
		cleaner: cleaner,
	}
}

func (c *Console) Run(cancel context.CancelFunc) {
	for {
		var command string
		var data string
		fmt.Printf("\033[38;2;76;121;72m%s\033[0m:\033[38;2;255;255;255m%s\033[0m$ ", c.fs.Name, c.fs.GetPath())
		fmt.Scanln(&command, &data)

		var err error
		switch command {
		case "cd":
			err = c.fs.Cd(data)
		case "mkdir":
			err = c.fs.Mkdir(data)
		case "tree":
			res, err := c.fs.Tree(data)
			if err == nil {
				fmt.Println(res)
			}
		case "ls":
			var res []string
			res, err = c.fs.Ls(data)
			if err == nil {
				fmt.Println("\033[38;2;189;38;93m", strings.Join(res, "\t"), "\033[0m")
			}
		case "rmdir":
			err = c.fs.Rmdir(data)
		case "mkf":
			err = c.fs.Makefile(data)
		case "rm":
			err = c.fs.Rm(data)
		case "of":
			file := new(os.File)
			file, err = c.fs.Of(data)
			defer file.Close()
			if err == nil {
				err = ide.Run(file)
				file.Close()
			}
		case "exit":
			cancel()
			return
		case "clear":
			c.cleaner.Clear()
		default:
			if command != "" {
				err = vanerrors.NewSimple("command not exist", command)
			}
		}
		if err != nil {
			fmt.Print("\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
		}
		c.fs.Save()
		c.fs.Reload()
		time.Sleep(time.Millisecond * 5)
	}
}
