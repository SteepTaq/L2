package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb_l2/18/internal/api/handler"
	"wb_l2/18/internal/api/middleware"
	"wb_l2/18/internal/config"
	"wb_l2/18/internal/repository"
	"wb_l2/18/internal/service"
)

type App struct {
	server *http.Server
	config *config.Config
}

func NewApp(config *config.Config) *App {
	repo := repository.NewRepository(repository.InMemory)

	service := service.NewService(repo)

	h := handler.NewHandler(service)
	handler.RegisterHandlers(h)

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: middleware.Chain(h.HTTPHandler(), middleware.Log),
		// TODO: ensure timeouts
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		server: server,
		config: config,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	errorChan := make(chan error, 1)
	go func() {
		slog.Info("Starting server...")
		if err := a.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			errorChan <- err
		}
	}()

	select {
	case err := <-errorChan:
		return fmt.Errorf("Unable to start server: %s", err)
	case <-time.After(time.Second):
	}

	slog.Info("Listening on http://localhost:" + a.config.Port)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	fmt.Println()
	slog.Info("Shutting down the server...")
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Unable to shutdown gracefully: %s", err)
	}

	return nil
}
