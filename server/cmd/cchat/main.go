package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sergey-frey/cchat/internal/app"
	"github.com/sergey-frey/cchat/internal/config"
	"github.com/sergey-frey/cchat/internal/http-server/handlers/session"
	"github.com/sergey-frey/cchat/internal/http-server/middleware/cors"
	"github.com/sergey-frey/cchat/internal/http-server/middleware/jwtcheck"
	authHandler "github.com/sergey-frey/cchat/internal/http-server/handlers/auth"
	userHandler "github.com/sergey-frey/cchat/internal/http-server/handlers/user"
	authService "github.com/sergey-frey/cchat/internal/services/auth"
	userService "github.com/sergey-frey/cchat/internal/services/user"
	"github.com/sergey-frey/cchat/internal/storage/postgres"
	"github.com/sergey-frey/cchat/cmd/migrator"
	"github.com/sergey-frey/cchat/internal/lib/logger/slogpretty"
)

// @title Cchat App API
// @version 0.1
// @description API Server for Cchat application

// @host localhost:8040
// @BasePath /cchat

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(cors.NewCORS)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	storagePath := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username, cfg.Storage.DBName, os.Getenv("PG_DB_PASSWORD"), cfg.Storage.SSLMode)

	pool, err := postgres.New(context.Background(), storagePath)
	if err != nil {
		panic(err)
	}
	authService := authService.New(pool, log)
	authHandler := authHandler.New(authService, log)

	userService := userService.New(pool, log)
	userHandler := userHandler.New(userService, log)

	migrator.NewMigration("postgres://postgres:qwerty@psql:5432/postgres?sslmode=disable", os.Getenv("MIGRATIONS_PATH"))

	// router.Get("/swagger/*", httpSwagger.Handler(
	// 	httpSwagger.URL("http://localhost:1323/swagger/doc.json"), //The url pointing to API definition
	// ))

	router.With(jwtcheck.JWTCheck).Route("/cchat/user", func(r chi.Router) {
		r.Get("/profile", userHandler.GetUser(context.Background()))
		r.Patch("/update", userHandler.UpdateUserInfo(context.Background()))
	})

	router.Route("/cchat/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login(context.Background()))
		r.Post("/register", authHandler.Register(context.Background()))
		r.Post("/session", session.CheckSession(context.Background(), log))
		r.Post("/logout", session.FinishSession(context.Background(), log))
	})

	log.Info("starting server")

	application := app.New(log, cfg, router)

	go func() {
		application.HTTPServer.Run()
	}()

	//Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop

	log.Info("stopping application", slog.String("signal", signal.String()))

	application.HTTPServer.Stop(context.Background())

	postgres.Close(context.Background(), pool)

	log.Info("application stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettyLogger()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
