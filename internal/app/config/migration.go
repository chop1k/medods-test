package config

import "flag"

// MigrationConfig holds options for locating and running database
// migrations.
type MigrationConfig struct {
	Path string
}

// RegisterMigrationFlags registers the migration flags on fs and returns a
// MigrationConfig that is populated once fs.Parse is called.
//
// Flag / env var mapping:
//
//	--migrations-path  MIGRATIONS_PATH  (default "migrations")
func RegisterMigrationFlags(fs *flag.FlagSet) *MigrationConfig {
	cfg := &MigrationConfig{}

	fs.StringVar(&cfg.Path, "migrations-path", envOrDefault("MIGRATIONS_PATH", "migrations"),
		"Filesystem path to the migration files (env: MIGRATIONS_PATH)")

	return cfg
}
