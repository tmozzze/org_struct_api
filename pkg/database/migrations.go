package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
	"github.com/tmozzze/org_struct_api/internal/config"
)

// RunMigrations - apply database migrations using goose
func RunMigrations(cfg config.Config, db *sql.DB) error {
	const op = "database.RunMigrations"

	if err := goose.SetDialect(cfg.DBDialect); err != nil {
		return fmt.Errorf("%s: failed to apply migrations: %w", op, err)
	}

	if err := goose.Up(db, cfg.MigrationsDir); err != nil {
		return fmt.Errorf("%s: failed to apply migrations: %w", op, err)
	}

	return nil
}
