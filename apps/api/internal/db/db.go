// Package db provides the Postgres connection used by the 415 API and
// applies embedded schema migrations on startup.
package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Config holds the connection parameters for the Postgres database.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ConfigFromEnv builds a Config from the environment, defaulting to values
// that match infra/db/docker-compose.yaml for local development.
func ConfigFromEnv() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5332"),
		User:     getEnv("DB_USER", "dba"),
		Password: getEnv("DB_PASSWORD", "dba"),
		Name:     getEnv("DB_NAME", "bank"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func (c Config) dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

// Open connects to Postgres and verifies the connection with a ping.
func Open(cfg Config) (*sql.DB, error) {
	conn, err := sql.Open("pgx", cfg.dsn())
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return conn, nil
}
