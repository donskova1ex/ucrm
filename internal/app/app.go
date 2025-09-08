package app

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"ucrm/logger"
)

type App struct {
	logger logger.Logger
	server *http.Server
}

func NewApp(log logger.Logger) *App {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Context(), "Handling from health", slog.String("method", "health"))
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &App{
		logger: log,
		server: server,
	}
}

func (a *App) Start(ctx context.Context) error {
	a.logger.Info(ctx, "Starting server")
	errChan := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error(ctx, "Error starting server", err)
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info(ctx, "Shutting down server")
	case err := <-errChan:
		return err
	}
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error(ctx, "Error shutting down server", err)
		return err
	}
	a.logger.Info(ctx, "Server shutdown complete")
	return nil
}
