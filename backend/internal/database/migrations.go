package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type MigrationManager struct {
	db *sql.DB
}

type Migration struct {
	ID       string
	Filename string
	SQL      string
}

type MigrationRepository struct {
	db *sql.DB
}

type MigrationFileLoader struct {
	migrationsDir string
}

func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{db: db}
}

func NewMigrationRepository(db *sql.DB) *MigrationRepository {
	return &MigrationRepository{db: db}
}

func NewMigrationFileLoader(dir string) *MigrationFileLoader {
	return &MigrationFileLoader{migrationsDir: dir}
}

func getMigrationsPath() string {
	paths := []string{
		"internal/database/migrations",
		"migrations",
		"../internal/database/migrations",
		"../../internal/database/migrations",
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return "internal/database/migrations"
}

func (mm *MigrationManager) RunMigrations() error {
	if mm.db == nil {
		return fmt.Errorf("database connection not established")
	}

	repo := NewMigrationRepository(mm.db)
	if err := repo.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	loader := NewMigrationFileLoader(getMigrationsPath())
	migrations, err := loader.LoadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	executed, err := repo.GetExecutedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	return mm.executePendingMigrations(migrations, executed, repo)
}

func (mm *MigrationManager) executePendingMigrations(migrations []Migration, executed map[string]bool, repo *MigrationRepository) error {
	for _, migration := range migrations {
		if !executed[migration.ID] {
			if os.Getenv("LOG_LEVEL") != "error" {
				log.Printf("Executing migration: %s", migration.Filename)
			}

			if err := mm.executeMigration(migration, repo); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", migration.ID, err)
			}

			if os.Getenv("LOG_LEVEL") != "error" {
				log.Printf("Migration completed: %s", migration.Filename)
			}
		}
	}
	return nil
}

func (mm *MigrationManager) executeMigration(migration Migration, repo *MigrationRepository) error {
	if _, err := mm.db.Exec(migration.SQL); err != nil {
		return err
	}
	return repo.RecordMigration(migration.ID, migration.Filename)
}

func (mfl *MigrationFileLoader) LoadMigrations() ([]Migration, error) {
	files, err := os.ReadDir(mfl.migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migration, err := mfl.loadMigrationFile(file.Name())
			if err != nil {
				return nil, err
			}
			migrations = append(migrations, migration)
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	return migrations, nil
}

func (mfl *MigrationFileLoader) loadMigrationFile(filename string) (Migration, error) {
	content, err := os.ReadFile(filepath.Join(mfl.migrationsDir, filename))
	if err != nil {
		return Migration{}, fmt.Errorf("failed to read migration file %s: %w", filename, err)
	}

	return Migration{
		ID:       strings.TrimSuffix(filename, ".sql"),
		Filename: filename,
		SQL:      string(content),
	}, nil
}

func (mr *MigrationRepository) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id TEXT PRIMARY KEY,
			filename TEXT NOT NULL,
			executed_at TIMESTAMPTZ DEFAULT now()
		);
	`
	_, err := mr.db.Exec(query)
	return err
}

func (mr *MigrationRepository) GetExecutedMigrations() (map[string]bool, error) {
	query := `SELECT id FROM migrations`
	rows, err := mr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		executed[id] = true
	}

	return executed, nil
}

func (mr *MigrationRepository) RecordMigration(id, filename string) error {
	query := `INSERT INTO migrations (id, filename) VALUES ($1, $2)`
	_, err := mr.db.Exec(query, id, filename)
	return err
}
