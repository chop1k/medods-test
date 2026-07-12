package config

import (
	"flag"
	"fmt"
	"net/url"
	"time"
)

// DatabaseConfig holds all Postgres connection/pool tuning options. Every
// field can be set via a CLI flag or its corresponding environment
// variable (CLI flags win when both are provided).
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DSN builds a standard "postgres://" connection string from the config.
// It is pure formatting - actually opening a connection is left to
// internal/database (not implemented yet, see its TODOs).
func (c DatabaseConfig) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   "/" + c.Name,
	}
	q := u.Query()
	q.Set("sslmode", c.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

// RegisterDatabaseFlags registers the Postgres connection flags on fs and
// returns a DatabaseConfig that is populated once fs.Parse is called.
//
// Flag / env var mapping:
//
//	--db-host               DB_HOST                (default "localhost")
//	--db-port               DB_PORT                (default 5432)
//	--db-user               DB_USER                (default "postgres")
//	--db-password           DB_PASSWORD             (default "")
//	--db-name               DB_NAME                (default "trackmytasks")
//	--db-sslmode            DB_SSLMODE              (default "disable")
//	--db-max-open-conns     DB_MAX_OPEN_CONNS       (default 10)
//	--db-max-idle-conns     DB_MAX_IDLE_CONNS       (default 5)
//	--db-conn-max-lifetime  DB_CONN_MAX_LIFETIME    (default 5m)
func RegisterDatabaseFlags(fs *flag.FlagSet) *DatabaseConfig {
	cfg := &DatabaseConfig{}

	fs.StringVar(&cfg.Host, "db-host", envOrDefault("DB_HOST", "localhost"),
		"Postgres host (env: DB_HOST)")
	fs.IntVar(&cfg.Port, "db-port", envIntOrDefault("DB_PORT", 5432),
		"Postgres port (env: DB_PORT)")
	fs.StringVar(&cfg.User, "db-user", envOrDefault("DB_USER", "postgres"),
		"Postgres user (env: DB_USER)")
	fs.StringVar(&cfg.Password, "db-password", envOrDefault("DB_PASSWORD", ""),
		"Postgres password (env: DB_PASSWORD)")
	fs.StringVar(&cfg.Name, "db-name", envOrDefault("DB_NAME", "trackmytasks"),
		"Postgres database name (env: DB_NAME)")
	fs.StringVar(&cfg.SSLMode, "db-sslmode", envOrDefault("DB_SSLMODE", "disable"),
		"Postgres sslmode, e.g. disable/require/verify-full (env: DB_SSLMODE)")
	fs.IntVar(&cfg.MaxOpenConns, "db-max-open-conns", envIntOrDefault("DB_MAX_OPEN_CONNS", 10),
		"Maximum number of open Postgres connections (env: DB_MAX_OPEN_CONNS)")
	fs.IntVar(&cfg.MaxIdleConns, "db-max-idle-conns", envIntOrDefault("DB_MAX_IDLE_CONNS", 5),
		"Maximum number of idle Postgres connections (env: DB_MAX_IDLE_CONNS)")
	fs.DurationVar(&cfg.ConnMaxLifetime, "db-conn-max-lifetime", envDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		"Maximum lifetime of a pooled Postgres connection (env: DB_CONN_MAX_LIFETIME)")

	return cfg
}
