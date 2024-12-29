package main

import (
	"context"
	"os/signal"
	"syscall"
	"vfs/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Run(ctx)

}
