// Package migrate will run database schema migrations for the application.
// It is intentionally left unimplemented - only the configuration surface
// (see internal/config.DatabaseConfig and MigrationConfig) is wired up for
// now.
package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/chop1k/medods-test/internal/config"
)

type FileInfo struct {
	ID      int
	Name    string
	Content string
	Path    string
}

// Run is a placeholder for applying pending migrations found under
// migrationsPath against the database described by dbCfg.
//
// TODO: wire in a migration tool (e.g. golang-migrate/migrate) and run all
// pending migrations from migrationsPath against dbCfg.DSN().
func Run(dbCfg config.DatabaseConfig, migrationsPath string) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.SSLMode)

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return err
	}

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

		err = runTransaction(tx, migration)

		if err != nil {
			return err
		}
	}

	return nil
}

func runTransaction(tx *sql.Tx, migration FileInfo) error {
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
	// Читаем все файлы в папке
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения папки: %v", err)
	}

	var fileInfos []FileInfo

	// Обрабатываем каждый файл
	for _, entry := range entries {
		// Пропускаем директории
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// Проверяем расширение .sql
		if !strings.HasSuffix(fileName, ".sql") {
			continue
		}

		// Убираем расширение .sql
		nameWithoutExt := strings.TrimSuffix(fileName, ".sql")

		// Разбиваем на ID и имя (ищем первый _)
		underscoreIndex := strings.Index(nameWithoutExt, "_")
		if underscoreIndex == -1 {
			continue // Пропускаем файлы без _
		}

		// Парсим ID
		idStr := nameWithoutExt[:underscoreIndex]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue // Пропускаем файлы с невалидным ID
		}

		// Получаем имя (все что после первого _)
		name := nameWithoutExt[underscoreIndex+1:]

		// Полный путь к файлу
		fullPath := filepath.Join(folderPath, fileName)

		// Читаем содержимое файла
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения файла %s: %v", fileName, err)
		}

		// Создаем структуру
		fileInfo := FileInfo{
			ID:      id,
			Name:    name,
			Content: string(content),
			Path:    fullPath,
		}

		fileInfos = append(fileInfos, fileInfo)
	}

	// Сортируем по ID
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ID < fileInfos[j].ID
	})

	return fileInfos, nil
}
