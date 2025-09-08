package main

import (
	"context"
	"os/signal"
	"syscall"
	"ucrm/internal/app"
	"ucrm/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logger.NewSlog()
	newApp := app.NewApp(log)

	if err := newApp.Start(ctx); err != nil {
		log.Error(ctx, "error starting app", err)
	}
}
