package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/vandi37/vanfs/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Run(ctx)

}
