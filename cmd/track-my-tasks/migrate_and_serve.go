package main

import (
	"flag"

	"github.com/chop1k/medods-test/internal/config"
	"github.com/chop1k/medods-test/internal/database"
)

// runMigrateAndServe implements the `migrate-and-serve` command: it applies
// pending migrations and then starts the HTTP server in the same process.
func runMigrateAndServe(args []string) {
	fs := flag.NewFlagSet("migrate-and-serve", flag.ExitOnError)
	serverCfg := config.RegisterServerFlags(fs)
	dbCfg := config.RegisterDatabaseFlags(fs)
	migCfg := config.RegisterMigrationFlags(fs)
	if err := fs.Parse(args); err != nil {
		fatal("migrate-and-serve: failed to parse flags: %v", err)
	}

	db, err := database.Connect(*dbCfg)

	if err != nil {
		fatal("serve: failed to connect to db: %v", err)
	}

	applyMigrations(*dbCfg, *migCfg)
	startServer(*serverCfg, db)
}
