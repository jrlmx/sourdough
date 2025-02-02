package main

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func ensureDBExists(dbPath string) error {
	// Create the sourdough directory if it doesn't exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dbPath, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// Create the sourdough database if it doesn't exist
	if _, err := os.Stat(filepath.Join(dbPath, "sourdough.db")); os.IsNotExist(err) {
		if _, err := os.Create(filepath.Join(dbPath, "sourdough.db")); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func ensureStarterTableExists(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS starters (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		description TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TRIGGER IF NOT EXISTS update_starters_updated_at
	AFTER UPDATE ON starters
	BEGIN
		UPDATE starters SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;
	`)
	if err != nil {
		return err
	}
	return nil
}

func database() (*sql.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(homeDir, ".sourdough")
	if err := ensureDBExists(dbPath); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(dbPath, "sourdough.db"))
	if err != nil {
		return nil, err
	}
	if err := ensureStarterTableExists(db); err != nil {
		return nil, err
	}
	return db, nil
}
