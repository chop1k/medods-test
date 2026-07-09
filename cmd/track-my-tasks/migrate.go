package main

import (
	"flag"
	"log"

	"github.com/chop1k/medods-test/internal/config"
	"github.com/chop1k/medods-test/internal/migrations"
)

// runMigrate implements the `migrate` command: it applies pending database
// migrations and exits.
func runMigrate(args []string) {
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	dbCfg := config.RegisterDatabaseFlags(fs)
	migCfg := config.RegisterMigrationFlags(fs)
	if err := fs.Parse(args); err != nil {
		fatal("migrate: failed to parse flags: %v", err)
	}

	applyMigrations(*dbCfg, *migCfg)
}

// applyMigrations runs pending migrations, exiting the process on failure.
func applyMigrations(dbCfg config.DatabaseConfig, migCfg config.MigrationConfig) {
	log.Printf("running migrations from %q against %s:%d/%s", migCfg.Path, dbCfg.Host, dbCfg.Port, dbCfg.Name)

	if err := migrations.Run(dbCfg, migCfg.Path); err != nil {
		fatal("migrate: %v", err)
	}

	log.Println("migrations applied successfully")
}
