package config

import (
	"flag"
	"fmt"
	"time"
)

// ServerConfig holds all HTTP server tuning options. Every field can be set
// via a CLI flag or its corresponding environment variable (CLI flags win
// when both are provided).
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Addr returns the "host:port" address the HTTP server should bind to.
func (c ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// RegisterServerFlags registers the HTTP server flags on fs and returns a
// ServerConfig that is populated once fs.Parse is called.
//
// Flag / env var mapping:
//
//	--http-host           HTTP_HOST           (default "0.0.0.0")
//	--http-port           HTTP_PORT           (default 8080)
//	--http-read-timeout   HTTP_READ_TIMEOUT   (default 15s)
//	--http-write-timeout  HTTP_WRITE_TIMEOUT  (default 15s)
//	--http-idle-timeout   HTTP_IDLE_TIMEOUT   (default 60s)
func RegisterServerFlags(fs *flag.FlagSet) *ServerConfig {
	cfg := &ServerConfig{}

	fs.StringVar(&cfg.Host, "http-host", envOrDefault("HTTP_HOST", "0.0.0.0"),
		"HTTP server bind host (env: HTTP_HOST)")
	fs.IntVar(&cfg.Port, "http-port", envIntOrDefault("HTTP_PORT", 8080),
		"HTTP server bind port (env: HTTP_PORT)")
	fs.DurationVar(&cfg.ReadTimeout, "http-read-timeout", envDurationOrDefault("HTTP_READ_TIMEOUT", 15*time.Second),
		"HTTP server read timeout (env: HTTP_READ_TIMEOUT)")
	fs.DurationVar(&cfg.WriteTimeout, "http-write-timeout", envDurationOrDefault("HTTP_WRITE_TIMEOUT", 15*time.Second),
		"HTTP server write timeout (env: HTTP_WRITE_TIMEOUT)")
	fs.DurationVar(&cfg.IdleTimeout, "http-idle-timeout", envDurationOrDefault("HTTP_IDLE_TIMEOUT", 60*time.Second),
		"HTTP server idle timeout (env: HTTP_IDLE_TIMEOUT)")

	return cfg
}
