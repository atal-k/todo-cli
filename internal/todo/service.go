// internal/todo/service.go
package todo

import (
	"fmt"
	"github.com/atalkumarme/todo-cli/internal/db"
)

type Service struct {
	db *db.DB
}

func NewService(db *db.DB) *Service {
	return &Service{db: db}
}

func (s *Service) AddTask(title string) error {
	_, err := s.db.CreateTask(title)
	return err
}

// ListAllTasks returns all tasks
func (s *Service) ListAllTasks() error {
	tasks, err := s.db.GetAllTasks()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		status := "[ ]"
		if task.Completed {
			status = "[✓]"
		}
		fmt.Printf("%s %d: %s (Created: %s)\n",
			status, task.ID, task.Title, task.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}

// ListByStatus returns tasks by their completion status
func (s *Service) ListByStatus(completed bool) error {
	tasks, err := s.db.GetTasksByStatus(completed)
	if err != nil {
		return err
	}
	status := "Pending"
	if completed {
		status = "Completed"
	}
	fmt.Printf("\n--- %s Tasks ---\n", status)
	for _, task := range tasks {
		fmt.Printf("%d: %s (Created: %s)\n",
			task.ID, task.Title, task.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}

// ListPendingTasks returns all incomplete tasks
func (s *Service) ListPendingTasks() error {
	tasks, err := s.db.GetTasksByStatus(false)
	if err != nil {
		return err
	}
	// Print tasks
	for _, task := range tasks {
		fmt.Printf("ID: %d | Title: %s | Created: %s\n",
			task.ID, task.Title, task.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}
