package main

import (
	"fmt"
	"log/slog"
	"os"
	_ "time/tzdata"

	"github.com/tmozzze/org_struct_api/internal/config"
	"github.com/tmozzze/org_struct_api/pkg/database"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

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

	fmt.Println("Vsyo ok")

	// Init Repos

	// Init Service

	// Start Server (net/http)
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
