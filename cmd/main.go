package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/vandi37/vanfs/internal/app"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Run(ctx)
	fmt.Print("\033[0m")
}
