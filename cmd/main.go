// cmd/todo/main.go
// Package main provides a CLI todo application
package main

import (
	"fmt"
	"github.com/atal-k/todo-cli/internal/db"
	"github.com/atal-k/todo-cli/internal/todo"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// getCommands returns all CLI commands
func getCommands(tm *todo.Service) []cli.Command {
	return []cli.Command{
		createAddCommand(tm),
		createListCommand(tm),
		createDoneCommand(tm),
		createDeleteCommand(tm),
		createClearCommand(tm),
	}
}

// createAddCommand creates a new task
func createAddCommand(tm *todo.Service) cli.Command {
	return cli.Command{
		Name:      "add",
		Aliases:   []string{"a", "new"},
		Usage:     "Add a new task to the list",
		ArgsUsage: "TASK_DESCRIPTION",
		Description: `Add a new task to your todo list.
			Example: todo add "Complete the project documentation"`,
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("Error: Task description required\nUsage: todo add \"your task here\"")
			}
			task := c.Args().First()
			if len(task) < 3 {
				return fmt.Errorf("Error: Task description too short (minimum 3 characters)")
			}
			if err := tm.AddTask(task); err != nil {
				return fmt.Errorf("Failed to add task: %v", err)
			}
			fmt.Printf("Successfully added task: %s\n", task)
			return nil
		},
	}
}

// createListCommand handles task listing with different filters
func createListCommand(tm *todo.Service) cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"l", "ls"},
		Usage:   "List all tasks",
		Description: `List tasks with optional filters:
			--all: Show all tasks
			--completed: Show only completed tasks
			--pending: Show only pending tasks
			Example: todo list --pending`,
		Flags: []cli.Flag{
			cli.BoolFlag{Name: flagAll, Usage: "Show all tasks"},
			cli.BoolFlag{Name: flagCompleted, Usage: "Show completed tasks"},
			cli.BoolFlag{Name: flagPending, Usage: "Show pending tasks"},
		},
		Action: func(c *cli.Context) error {
			switch {
			case c.Bool(flagCompleted):
				return tm.ListByStatus(true)
			case c.Bool(flagPending):
				return tm.ListPendingTasks()
			default:
				return tm.ListAllTasks()
			}
		},
	}
}

// createDoneCommand marks a task as completed
func createDoneCommand(tm *todo.Service) cli.Command {
	return cli.Command{
		Name:      "done",
		Aliases:   []string{"d", "complete"},
		Usage:     "Mark a task as completed",
		ArgsUsage: "TASK_ID",
		Description: `Mark a specific task as completed using its ID.
			Example: todo done 1`,
		Action: func(c *cli.Context) error {
			id, err := parseID(c)
			if err != nil {
				return fmt.Errorf("Invalid task ID: %v\nUsage: todo done TASK_ID", err)
			}
			if err := tm.CompleteTask(id); err != nil {
				return fmt.Errorf("Failed to complete task: %v", err)
			}
			fmt.Printf("Task %d marked as completed\n", id)
			return nil
		},
	}
}

// createDeleteCommand removes a task
func createDeleteCommand(tm *todo.Service) cli.Command {
	return cli.Command{
		Name:      "delete",
		Aliases:   []string{"rm", "remove"},
		Usage:     "Delete a task",
		ArgsUsage: "TASK_ID",
		Description: `Delete a specific task using its ID.
			Example: todo delete 1`,
		Action: func(c *cli.Context) error {
			id, err := parseID(c)
			if err != nil {
				return fmt.Errorf("Invalid task ID: %v\nUsage: todo delete TASK_ID", err)
			}
			if err := tm.DeleteTask(id); err != nil {
				return fmt.Errorf("Failed to delete task: %v", err)
			}
			fmt.Printf("Task %d deleted successfully\n", id)
			return nil
		},
	}
}

// createClearCommand removes all tasks
func createClearCommand(tm *todo.Service) cli.Command {
	return cli.Command{
		Name:  "clear",
		Usage: "Delete all tasks",
		Description: `Remove all tasks from the todo list.
			Warning: This action cannot be undone.
			Example: todo clear`,
		Action: func(c *cli.Context) error {
			fmt.Print("Are you sure you want to delete all tasks? [y/N]: ")
			var response string
			fmt.Scanln(&response)
			if response = strings.ToLower(response); response == "y" || response == "yes" {
				if err := tm.ClearTasks(); err != nil {
					return fmt.Errorf("Failed to clear tasks: %v", err)
				}
				fmt.Println("All tasks have been deleted")
				return nil
			}
			fmt.Println("Operation cancelled")
			return nil
		},
	}
}

// parseID validates and converts string ID to integer
func parseID(c *cli.Context) (int, error) {
	if c.NArg() == 0 {
		return 0, fmt.Errorf("task ID is required")
	}
	id, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return 0, fmt.Errorf("invalid task ID format")
	}
	if id < 1 {
		return 0, fmt.Errorf("task ID must be positive")
	}
	return id, nil
}
