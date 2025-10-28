package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sergey-frey/cchat/server/chat-service/internal/config"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
}

func New(log *slog.Logger, cfg *config.Config, router chi.Router) *App {
	httpServer := http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.Idle_timeout,
	}

	return &App{log: log, httpServer: &httpServer}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.httpServer.Addr),
	)

	log.Info("starting http server")

	if err := a.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("http server started")

	return nil
}

func (a *App) Stop(ctx context.Context) {
	const op = "httpapp.Stop"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("stoping http server")

	a.httpServer.Shutdown(ctx)
}
