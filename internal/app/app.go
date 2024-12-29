package app

import (
	"context"
	"fmt"
	"vfs/internal/console"
	"vfs/pkg/cleaner"
	"vfs/pkg/filesystem"
)

func Run(ctx context.Context) {
	cleaner := cleaner.New()
	cleaner.Clear()
	f, _ := filesystem.New("data.json")
	defer f.Source.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	console := console.New(f, cleaner)

	go console.Run(cancel)

	<-ctx.Done()
	fmt.Println("\n\033[38;2;196;124;25mExiting...\033[0m")
}
