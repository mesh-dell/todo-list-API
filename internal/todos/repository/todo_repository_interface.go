package repository

import (
	"context"

	"github.com/mesh-dell/todo-list-API/internal/todos"
	"github.com/mesh-dell/todo-list-API/internal/utils"
)

type ITodoRepository interface {
	Create(ctx context.Context, todoItem *todos.TodoItem) error
	FindByID(ctx context.Context, id uint) (todos.TodoItem, error)
	FindAllForUser(ctx context.Context, userID uint, pagination *utils.Pagination) (*utils.Pagination, error)
	Update(ctx context.Context, id uint, todoItem todos.TodoItem) error
	Delete(ctx context.Context, id uint) error
}
