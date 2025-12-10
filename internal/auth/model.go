package auth

import (
	"github.com/mesh-dell/todo-list-API/internal/todos"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"unique"`
	PasswordHash string
	TodoItems    []todos.TodoItem
}
