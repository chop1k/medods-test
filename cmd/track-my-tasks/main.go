package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/chop1k/medods-test/internal/app"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd, args := os.Args[1], os.Args[2:]

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	switch cmd {
	case "serve":
		app.Serve(args)
	case "seed":
		app.Seed(args)
	case "migrate":
		app.Migrate(args)
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

Run 'api <command> -h' for a full list of flags/env vars for that command.
`)
}
