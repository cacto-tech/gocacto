package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cacto-cms/app/infrastructure/database"
	"cacto-cms/app/infrastructure/database/seeds"
	"cacto-cms/config"
)

func main() {
	// Check for "migrate:fresh" command
	var migrateFresh bool
	var seed bool

	for _, arg := range os.Args[1:] {
		if arg == "migrate:fresh" {
			migrateFresh = true
		}
		if arg == "--seed" {
			seed = true
		}
	}

	if !migrateFresh {
		fmt.Println("Usage: artisan migrate:fresh [--seed]")
		fmt.Println("  migrate:fresh  - Drop all tables and re-run all migrations")
		fmt.Println("  --seed        - Run seeders after migration")
		os.Exit(1)
	}

	// Load config
	cfg := config.Load()

	fmt.Println(`
╔═════════════════════════════════════════╗
║         Cacto CMS Artisan CLI           ║
╚═════════════════════════════════════════╝
	`)

	// Get database path
	dbPath := cfg.DBPath
	if dbPath == "" {
		dbPath = "./cacto.db"
	}

	// If migrate:fresh, remove existing database file
	if migrateFresh {
		if _, err := os.Stat(dbPath); err == nil {
			log.Printf("🗑️  Removing existing database: %s", dbPath)
			if err := os.Remove(dbPath); err != nil {
				log.Fatalf("Failed to remove database: %v", err)
			}

			// Also remove WAL and SHM files
			walPath := dbPath + "-wal"
			shmPath := dbPath + "-shm"
			os.Remove(walPath)
			os.Remove(shmPath)
		}
	}

	// Create database and run migrations
	log.Printf("📂 Database path: %s", dbPath)

	// Ensure database directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Migrations completed")

	// Run seeders if requested
	if seed {
		log.Println("\n🌱 Running seeders...")
		seeder := seeds.NewSeeder(db.DB)
		if err := seeder.SeedAll(); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		log.Println("✅ Seeding completed")
	}

	fmt.Println("\n✅ Done!")
}
