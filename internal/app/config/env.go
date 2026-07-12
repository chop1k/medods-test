package config

import (
	"os"
	"strconv"
	"time"
)

// envOrDefault returns the value of the given environment variable, or def
// if it is unset/empty. It is used as the *default* value for a flag, so
// command-line arguments always take precedence over environment variables,
// which in turn take precedence over the hardcoded default.
func envOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func envIntOrDefault(key string, def int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func envDurationOrDefault(key string, def time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
