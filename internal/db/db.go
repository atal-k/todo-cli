// internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"time"
)

type Task struct {
	ID        int64
	Title     string
	Completed bool
	CreatedAt time.Time
}

type DB struct {
	db *sql.DB
}

// NewDB creates a new DB instance with proper configuration
func NewDB(dbPath string) (*DB, error) {
	// Open database with additional parameters for better concurrent access
	db, err := sql.Open("sqlite3", dbPath+"?_journal=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(1)            // SQLite supports only one writer
	db.SetMaxIdleConns(1)            // Keep connection alive
	db.SetConnMaxLifetime(time.Hour) // Recreate connections every hour

	// Verify database connection
	if err := db.Ping(); err != nil {
		db.Close() // Clean up before return
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{db: db}, nil
}

// Initialize creates the tasks table if it doesn't exist
func (d *DB) Initialize() error {
	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// Core database operations

func (d *DB) CreateTask(title string) (*Task, error) {
	res, err := d.db.Exec("INSERT INTO tasks (title) VALUES (?)", title)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Task{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}, nil
}
