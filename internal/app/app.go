package app

import (
	"context"
	"fmt"
	"vfs/internal/console"
	"vfs/internal/path"
	"vfs/pkg/cleaner"
	"vfs/pkg/filesystem"
)

func Run(ctx context.Context) {
	path, err := path.LoadPath()
	if err != nil {
		fmt.Print("\033[48;2;120;24;0;38;2;255;221;212m", err, "\033[0m\n")
		return
	}
	cleaner := cleaner.New()
	cleaner.Clear()
	f, _ := filesystem.New(path)
	defer f.Source.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	console := console.New(f, cleaner)

	go console.Run(cancel)

	<-ctx.Done()
	fmt.Println("\n\033[38;2;196;124;25mExiting...\033[0m")
}
