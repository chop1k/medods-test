package app

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/chop1k/medods-test/internal/app/config"
	"github.com/chop1k/medods-test/internal/transport/http"
)

func Migrate(args []string) {
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)

	cfg := config.RegisterConfigFlags(fs)

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	if err := migrate(cfg.DB, cfg.Migrations.Path); err != nil {
		panic(err)
	}
}

func Seed(args []string) {
	fs := flag.NewFlagSet("seed", flag.ExitOnError)

	cfg := config.RegisterConfigFlags(fs)

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	if err := seed(cfg.DB, cfg.Seeds.Path); err != nil {
		panic(err)
	}
}

func Serve(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)

	cfg := config.RegisterConfigFlags(fs)

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	db, err := OpenDatabase(cfg.DB)

	if err != nil {
		panic(err)
	}

	srv := http.NewServer(cfg.HTTP, db)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Listen(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()

	if err = srv.Shutdown(); err != nil {
		panic(err)
	}
}

func OpenDatabase(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	conn, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = conn.Ping()

	if err != nil {
		return nil, err
	}

	return conn, nil
}

type FileInfo struct {
	ID      int
	Name    string
	Content string
	Path    string
}

func seed(cfg *config.DatabaseConfig, seedsPath string) error {
	db, err := OpenDatabase(cfg)

	if err != nil {
		return err
	}
	defer db.Close()

	seeds, err := parseAndReadFiles(seedsPath)

	if err != nil {
		return err
	}

	for _, seed := range seeds {
		tx, err := db.Begin()

		if err != nil {
			return err
		}

		var id int

		err = db.QueryRow("select id from \"app\".\"seeds\" where id = $1", seed.ID).Scan(&id)

		if err == nil {
			continue
		}

		err = runSeedTransaction(tx, seed)

		if err != nil {
			return errors.Join(fmt.Errorf("got error at seed %d", seed.ID), err)
		}
	}

	return nil
}

func runSeedTransaction(tx *sql.Tx, seed FileInfo) error {
	defer tx.Rollback()

	_, err := tx.Exec(seed.Content)

	if err != nil {
		return err
	}

	_, err = tx.Exec("insert into \"app\".\"seeds\" (\"id\", \"name\", \"created_at\") values ($1, $2, now())", seed.ID, seed.Name)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func migrate(cfg *config.DatabaseConfig, migrationsPath string) error {
	db, err := OpenDatabase(cfg)

	if err != nil {
		return err
	}
	defer db.Close()

	migrations, err := parseAndReadFiles(migrationsPath)

	if err != nil {
		return err
	}

	for _, migration := range migrations {
		tx, err := db.Begin()

		if err != nil {
			return err
		}

		var id int

		err = db.QueryRow("select id from \"app\".\"migrations\" where id = $1", migration.ID).Scan(&id)

		if err == nil {
			continue
		}

		err = runMigrationTransaction(tx, migration)

		if err != nil {
			return errors.Join(fmt.Errorf("got error at migration %d", migration.ID), err)
		}
	}

	return nil
}

func runMigrationTransaction(tx *sql.Tx, migration FileInfo) error {
	defer tx.Rollback()

	_, err := tx.Exec(migration.Content)

	if err != nil {
		return err
	}

	_, err = tx.Exec("insert into \"app\".\"migrations\" (\"id\", \"name\", \"created_at\") values ($1, $2, now())", migration.ID, migration.Name)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func parseAndReadFiles(folderPath string) ([]FileInfo, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	var fileInfos []FileInfo

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		if !strings.HasSuffix(fileName, ".sql") {
			continue
		}

		nameWithoutExt := strings.TrimSuffix(fileName, ".sql")

		underscoreIndex := strings.Index(nameWithoutExt, "_")
		if underscoreIndex == -1 {
			continue
		}

		idStr := nameWithoutExt[:underscoreIndex]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}

		name := nameWithoutExt[underscoreIndex+1:]

		fullPath := filepath.Join(folderPath, fileName)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %v", fileName, err)
		}

		fileInfo := FileInfo{
			ID:      id,
			Name:    name,
			Content: string(content),
			Path:    fullPath,
		}

		fileInfos = append(fileInfos, fileInfo)
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ID < fileInfos[j].ID
	})

	return fileInfos, nil
}
