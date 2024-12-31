package app

import (
	"context"
	"fmt"
	"strings"
	"vfs/internal/console"
	"vfs/internal/path"
	"vfs/pkg/choose"
	"vfs/pkg/cleaner"
	"vfs/pkg/filesystem"
	"vfs/pkg/init_system"

	"github.com/AlecAivazis/survey/v2"
)

func Run(ctx context.Context) {

	variant, err := choose.Choose()
	if err != nil {
		fmt.Print("\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
		return
	}
	var fs = new(filesystem.Filesystem)

	switch variant {
	case 0:
		path, err := path.LoadPath()

		if err != nil {
			fmt.Print("\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
			return
		}
		fs, err = filesystem.New(path)

		if err != nil {
			fmt.Print("\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
			return
		}
	case 1:
		prompt := &survey.Input{
			Message: "Backup path:",
		}
		var path string
		survey.AskOne(prompt, &path)

		if !strings.HasSuffix(path, "/") {
			path += "/"
		}

		fs, err = filesystem.New(path)

		if err != nil {
			fmt.Print("\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
			return
		}
	case 2:
		fs, err = init_system.Init()
		if err != nil {
			fmt.Print("\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
			return
		}
	default:
		fmt.Println("\033[38;2;196;124;25mExiting...")
		return
	}

	cleaner := cleaner.New()
	cleaner.Clear()

	console := console.New(fs, cleaner)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go console.Run(cancel)

	<-ctx.Done()
	fmt.Println("\n\033[38;2;196;124;25mExiting...\033[0m")
}
