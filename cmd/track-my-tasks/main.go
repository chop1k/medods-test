// Command api is the entrypoint for the Track My Tasks HTTP API.
//
// Usage:
//
//	api serve              Run the HTTP server
//	api migrate            Run pending database migrations
//	api migrate-and-serve  Run migrations, then start the HTTP server
//
// Every option accepts both a CLI flag and an environment variable; run
// `api <command> -h` to see the full list for that command.
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd, args := os.Args[1], os.Args[2:]

	switch cmd {
	case "serve":
		runServe(args)
	case "migrate":
		runMigrate(args)
	case "migrate-and-serve":
		runMigrateAndServe(args)
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprint(os.Stderr, `Track My Tasks API

Usage:
  api <command> [flags]

Commands:
  serve              Run the HTTP server
  migrate            Run pending database migrations
  migrate-and-serve  Run migrations, then start the HTTP server

Run 'api <command> -h' for a full list of flags/env vars for that command.
`)
}

// fatal logs a fatal error and exits. Kept as a single choke point so
// commands report failures consistently.
func fatal(format string, args ...any) {
	log.Fatalf(format, args...)
}
