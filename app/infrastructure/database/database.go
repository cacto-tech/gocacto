package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Database wraps the SQL database connection
type Database struct {
	DB *sql.DB
}

// New creates a new database connection and runs migrations
func New(dbPath string) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Set pragmas for better performance
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	database := &Database{DB: db}

	// Run migrations
	if err := database.runMigrations(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return database, nil
}

func (d *Database) runMigrations() error {
	log.Println("🔄 Running migrations...")

	files, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	// Sort files to ensure order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		content, err := migrationFS.ReadFile("migrations/" + file.Name())
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file.Name(), err)
		}

		if _, err := d.DB.Exec(string(content)); err != nil {
			return fmt.Errorf("migration %s failed: %w", file.Name(), err)
		}

		log.Printf("  ✓ Applied: %s", file.Name())
	}

	log.Println("✅ All migrations completed")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}
