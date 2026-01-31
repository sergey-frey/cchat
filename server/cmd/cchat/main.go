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
	// "github.com/sergey-frey/cchat/cmd/migrator"
	_ "github.com/sergey-frey/cchat/docs"
	"github.com/sergey-frey/cchat/internal/app"
	"github.com/sergey-frey/cchat/internal/config"
	authHandler "github.com/sergey-frey/cchat/internal/http-server/handlers/auth"
	chatHandler "github.com/sergey-frey/cchat/internal/http-server/handlers/chat"
	"github.com/sergey-frey/cchat/internal/http-server/handlers/session"
	userHandler "github.com/sergey-frey/cchat/internal/http-server/handlers/user"
	"github.com/sergey-frey/cchat/internal/http-server/middleware/cors"
	"github.com/sergey-frey/cchat/internal/http-server/middleware/jwtcheck"
	"github.com/sergey-frey/cchat/internal/lib/logger/slogpretty"
	authService "github.com/sergey-frey/cchat/internal/services/auth"
	chatService "github.com/sergey-frey/cchat/internal/services/chat"
	userService "github.com/sergey-frey/cchat/internal/services/user"
	"github.com/sergey-frey/cchat/internal/storage/postgres"
	"github.com/sergey-frey/cchat/internal/storage/redis"
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

	storagePath := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.PostgreStorage.Host, cfg.PostgreStorage.Port, cfg.PostgreStorage.Username, cfg.PostgreStorage.DBName, os.Getenv("PG_DB_PASSWORD"), cfg.PostgreStorage.SSLMode)

	redisPool, err := redis.New(context.Background(), cfg.RedisStorage)
	if err != nil {
		panic(err)
	}

	pool, err := postgres.New(context.Background(), storagePath)
	if err != nil {
		panic(err)
	}
	authService := authService.New(pool, log)
	authHandler := authHandler.New(authService, log)

	userService := userService.New(pool, log)
	userHandler := userHandler.New(userService, log)

	chatService := chatService.New(pool, redisPool, log)
	chatHandler := chatHandler.New(chatService, redisPool, log)

	// migrator.NewMigration("postgres://user:password@db:5432/mydb?sslmode=disable", os.Getenv("MIGRATIONS_PATH"))
	
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8040/swagger/doc.json"), //The url pointing to API definition
	))

	router.Route("/cchat/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login(context.Background()))
		r.Post("/register", authHandler.Register(context.Background()))
		r.Post("/session", session.CheckSession(context.Background(), log))
		r.Post("/logout", session.FinishSession(context.Background(), log))
	})

	router.With(jwtcheck.JWTCheck).Route("/cchat/user", func(r chi.Router) {
		r.Get("/myprofile", userHandler.MyProfile(context.Background()))
		r.Get("/profile/{username}", userHandler.Profile(context.Background()))
		r.Get("/list-profiles", userHandler.ListProfiles(context.Background()))
		r.Patch("/update", userHandler.UpdateInfo(context.Background()))
	})

	router.With(jwtcheck.JWTCheck).Route("/cchat/chat", func(r chi.Router) {
		r.Post("/new", chatHandler.NewChat(context.Background()))
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
