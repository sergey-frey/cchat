package main

import (
	"context"
	"fmt"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sergey-frey/cchat/user-service/cmd/migrator"
	_ "github.com/sergey-frey/cchat/user-service/docs"
	"github.com/sergey-frey/cchat/user-service/internal/app"
	"github.com/sergey-frey/cchat/user-service/internal/config"
	"github.com/sergey-frey/cchat/user-service/internal/http-server/handlers/system"
	userHandler "github.com/sergey-frey/cchat/user-service/internal/http-server/handlers/user"
	"github.com/sergey-frey/cchat/user-service/internal/http-server/middleware/cors"
	"github.com/sergey-frey/cchat/user-service/internal/http-server/middleware/jwtcheck"
	"github.com/sergey-frey/cchat/user-service/internal/lib/logger/slogpretty"
	"github.com/sergey-frey/cchat/user-service/internal/provider/storage/postgres"
	userService "github.com/sergey-frey/cchat/user-service/internal/services/user"
	"github.com/swaggo/http-swagger/v2"
)

// @title Cchat App API
// @version 0.1
// @description API Server for Cchat application

// @host localhost:8040
// @BasePath /cchat

// @securityDefinitions.cookie CookieAuth
// @in cookie
// @name accessToken

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

	storagePath := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), "disable")

	pool, err := postgres.New(context.Background(), storagePath)
	if err != nil {
		panic(err)
	}

	userService := userService.New(pool, log)
	userHandler := userHandler.New(userService, log)

	migrator.NewMigration("postgres://user:password@users-db:5432/usersdb?sslmode=disable", os.Getenv("MIGRATIONS_PATH"))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8040/swagger/doc.json"), //The url pointing to API definition
	))

	router.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser(context.Background()))
		r.Get("/{uuid}", userHandler.GetUserByID(context.Background()))
		r.Get("/", userHandler.GetUserByEmail(context.Background()))
	})

	router.With(jwtcheck.JWTCheck).Route("/profiles", func(r chi.Router) {
		r.Get("/", userHandler.Profiles(context.Background()))
		// r.Patch("/{id}", userHandler.UpdateInfo(context.Background()))
	})

	router.Get("/system-id", system.GetSystemID(context.Background(), log))

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
