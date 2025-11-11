package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sergey-frey/cchat/server/auth-service/cmd/migrator"
	_ "github.com/sergey-frey/cchat/server/auth-service/docs"
	"github.com/sergey-frey/cchat/server/auth-service/internal/app"
	"github.com/sergey-frey/cchat/server/auth-service/internal/config"
	authHandler "github.com/sergey-frey/cchat/server/auth-service/internal/http-server/handlers/auth"
	"github.com/sergey-frey/cchat/server/auth-service/internal/http-server/handlers/session"
	"github.com/sergey-frey/cchat/server/auth-service/internal/http-server/middleware/cors"
	"github.com/sergey-frey/cchat/server/auth-service/internal/lib/logger/slogpretty"
	"github.com/sergey-frey/cchat/server/auth-service/internal/provider/api/userapi"
	"github.com/sergey-frey/cchat/server/auth-service/internal/provider/storage/postgres"
	authService "github.com/sergey-frey/cchat/server/auth-service/internal/services/auth"
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

	apiHttpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	userApiClient := userapi.NewClient(apiHttpClient, os.Getenv("USERS_SERVICE_1_URL"), log)

	authService := authService.New(pool, userApiClient, log)
	authHandler := authHandler.New(authService, log)

	migrator.NewMigration("postgres://user:password@auth-db:5432/authdb?sslmode=disable", os.Getenv("MIGRATIONS_PATH"))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8040/swagger/doc.json"), //The url pointing to API definition
	))

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login(context.Background()))
		r.Post("/register", authHandler.Register(context.Background()))
		r.Post("/session", session.CheckSession(context.Background(), log))
		r.Post("/logout", session.FinishSession(context.Background(), log))
		r.Patch("/password/change", authHandler.ChangePassword(context.Background()))
		r.Post("/password/reset", authHandler.ResetPassword(context.Background()))
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
