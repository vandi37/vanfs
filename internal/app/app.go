package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/vandi37/vanfs/internal/console"
	"github.com/vandi37/vanfs/internal/load"
	"github.com/vandi37/vanfs/internal/path"
	"github.com/vandi37/vanfs/pkg/choose"
	"github.com/vandi37/vanfs/pkg/cleaner"
	"github.com/vandi37/vanfs/pkg/filesystem"
	"github.com/vandi37/vanfs/pkg/init_system"
)

func Run(ctx context.Context) {
	variant, ok, err := load.Load()
	if err != nil {
		fmt.Fprint(os.Stderr, "\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
	}
	if ok {
		return
	}
	if variant < 0 {
		variant, err = choose.Choose()
		if err != nil {
			fmt.Fprint(os.Stderr, "\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
			return
		}
	}

	var fs = new(filesystem.Filesystem)

	switch variant {
	case 0:
		path, err := path.LoadPath()

		if err != nil {
			fmt.Fprint(os.Stderr, "\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
			return
		}
		fs, err = filesystem.New(path)

		if err != nil {
			fmt.Fprint(os.Stderr, "\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
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
			fmt.Fprint(os.Stderr, "\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
			return
		}
	case 2:
		fs, err = init_system.Init()
		if err != nil {
			fmt.Fprint(os.Stderr, "\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
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
	fs.Source.Close()
	fmt.Println("\n\033[38;2;196;124;25mExiting...\033[0m")
}
