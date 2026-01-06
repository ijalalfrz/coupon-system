package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ijalalfrz/coupon-system/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

// DB wraps sqlx.DB with additional functionality
type DB struct {
	*sqlx.DB
}

// InitDB initializes PostgreSQL database connection with sqlx.
func InitDB(cfg *config.Config) (*DB, error) {
	slog.Debug("connecting to DB...", slog.String("dsn", cfg.DB.DSN))

	db, err := sqlx.Connect("pgx", cfg.DB.DSN)
	if err != nil {
		slog.Error("could not connect to DB", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.DB.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.DB.MaxConnectionLifetime)
	db.SetConnMaxIdleTime(cfg.DB.MaxConnectionIdleTime)

	slog.Info("database connection pool configured",
		slog.Int("max_open_conns", cfg.DB.MaxOpenConnections),
		slog.Int("max_idle_conns", cfg.DB.MaxIdleConnections),
		slog.Duration("conn_max_lifetime", cfg.DB.MaxConnectionLifetime),
		slog.Duration("conn_max_idle_time", cfg.DB.MaxConnectionIdleTime),
	)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("could not ping DB", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("successfully connected to database")

	return &DB{DB: db}, nil
}

// HealthCheck performs a database health check
func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Stats returns database connection pool statistics
func (db *DB) Stats() map[string]interface{} {
	stats := db.DB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}
