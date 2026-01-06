package config

import (
	"log/slog"
	"time"
)

type LogLeveler string

func (l LogLeveler) Level() slog.Level {
	var level slog.Level

	_ = level.UnmarshalText([]byte(l))

	return level
}

// Config holds the server configuration.
type Config struct {
	LogLevel LogLeveler `mapstructure:"LOG_LEVEL"`
	DB       DB         `mapstructure:",squash"`
	HTTP     HTTP       `mapstructure:",squash"`
}
type DB struct {
	DSN                   string        `mapstructure:"DB_DSN"`
	MaxOpenConnections    int           `mapstructure:"DB_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections    int           `mapstructure:"DB_MAX_IDLE_CONNECTIONS"`
	MaxConnectionLifetime time.Duration `mapstructure:"DB_MAX_CONNECTIONS_LIFETIME"`
	MaxConnectionIdleTime time.Duration `mapstructure:"DB_MAX_CONNECTION_IDLE_TIME"`
}

type HTTP struct {
	Port    int           `mapstructure:"HTTP_PORT"`
	Timeout time.Duration `mapstructure:"HTTP_TIMEOUT"`
}
