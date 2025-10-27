package main

import (
	"context"
	"log/slog"
	"os"
	"wb_l2/18/internal/app"
	"wb_l2/18/internal/config"
)

func main() {
	config, err := config.ReadConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	app := app.NewApp(config)

	if err := app.Run(context.Background()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
