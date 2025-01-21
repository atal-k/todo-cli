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
