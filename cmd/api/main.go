package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/go-playground/validator/v10"
	"github.com/tmozzze/org_struct_api/internal/config"
	httpHandler "github.com/tmozzze/org_struct_api/internal/handler/http"
	"github.com/tmozzze/org_struct_api/internal/repository/postgres"
	"github.com/tmozzze/org_struct_api/internal/service"
	"github.com/tmozzze/org_struct_api/pkg/database"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Swagger UI
// http://localhost:8080/swagger/index.html

// @title Organization Structure API
// @version 1.0
// @description API server for Organization Structure
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {

	// Init Config
	cfg := config.MustLoad()

	// Init logger (slog)
	log := setupLogger(cfg.Env)
	log.Info("starting org_struct_api", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init DB (GORM)
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Error("failed to init database", slog.Any("err", err))
		os.Exit(1)
	}
	log.Info("database is initialized")

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("failed to get sql.DB from gorm.DB", slog.Any("err", err))
		os.Exit(1)
	}
	// Close DB connection on exit
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Error("failed to close database connection", slog.Any("err", err))
		} else {
			log.Info("database connection closed")
		}
	}()

	// Run Migrations
	if err := database.RunMigrations(*cfg, sqlDB); err != nil {
		log.Error("failed to run migrations", slog.Any("err", err))
		os.Exit(1)
	}
	log.Info("migrations applied successfully")

	// Init Repos
	repo := postgres.NewRepository(db)

	// Init Validator
	validate := validator.New()

	// Init Service
	svc := service.NewService(repo, log, validate)

	// Init Handlers
	handler := httpHandler.NewHandler(svc, log)

	// Router
	router := httpHandler.NewRouter(handler)

	// Start Server
	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	log.Info("server started", slog.String("addr", cfg.HTTPServer.Address))

	// signal
	<-done
	log.Info("stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.String("err", err.Error()))
	}
}

func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal: // Text Debug
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev: // JSON Debug
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd: // JSON Info
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
