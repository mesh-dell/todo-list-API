package auth

import (
	"time"

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

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	Token     string `gorm:"uniqueIndex;type:varchar(500)"`
	ExpiresAt time.Time
}
