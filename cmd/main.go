// cmd/todo/main.go
// Package main provides a CLI todo application
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atalkumarme/todo-cli/internal/db"
	"github.com/atalkumarme/todo-cli/internal/todo"
	"github.com/urfave/cli"
)

// Constants for command-line flags
const (
	flagAll       = "all"
	flagCompleted = "completed"
	flagPending   = "pending"
)

func main() {

	// Initialize application
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get home directory: ", err)
	}

	dbPath := filepath.Join(homeDir, ".todo.db")
	// Setup database connection
	db, err := db.NewDB(dbPath)
	if err != nil {
		log.Fatal("Database initialization failed: ", err)
	}
	defer db.Close()

	if err := db.Initialize(); err != nil {
		log.Fatal("Failed to create tables: ", err)
	}

	// Initialize service layer
	todoManager := todo.NewService(db)

	// Configure CLI application
	app := cli.NewApp()
	app.Name = "todo"
	app.Usage = "A simple and efficient CLI todo manager"
	app.Version = "1.0.0"
	app.Authors = []cli.Author{
		{Name: "Todo CLI Team"},
	}

	// Register commands
	app.Commands = getCommands(todoManager)

	// Start CLI application
	if err := app.Run(os.Args); err != nil {
		log.Fatal("Application error: ", err)
	}
}
