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

func (d *DB) GetTask(id int64) (*Task, error) {
	row := d.db.QueryRow("SELECT id, title, completed, created_at FROM tasks WHERE id = ?", id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (d *DB) GetAllTasks() ([]*Task, error) {
	rows, err := d.db.Query("SELECT id, title, completed, created_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (d *DB) GetTasksByStatus(completed bool) ([]*Task, error) {
	rows, err := d.db.Query("SELECT id, title, completed, created_at FROM tasks WHERE completed = ?", completed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (d *DB) UpdateTaskStatus(id int64, completed bool) error {
	_, err := d.db.Exec("UPDATE tasks SET completed = ? WHERE id = ?", completed, id)
	return err
}

func (d *DB) UpdateTask(task *Task) error {
	_, err := d.db.Exec("UPDATE tasks SET title = ?, completed = ? WHERE id = ?", task.Title, task.Completed, task.ID)
	return err
}

func (d *DB) DeleteTask(id int64) error {
	_, err := d.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (d *DB) DeleteAllTasks() error {
	_, err := d.db.Exec("DELETE FROM tasks")
	return err
}
