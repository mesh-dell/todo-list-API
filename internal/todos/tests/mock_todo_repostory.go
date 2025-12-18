package tests

import (
	"context"

	"github.com/mesh-dell/todo-list-API/internal/todos"
	"github.com/mesh-dell/todo-list-API/internal/todos/dtos"
)

type MockTodoRepository struct {
	createFn   func(ctx context.Context, todoItem *todos.TodoItem) error
	findByIdFn func(ctx context.Context, id uint) (todos.TodoItem, error)
	updateFn   func(ctx context.Context, id uint, todoItem todos.TodoItem) error
	deleteFn   func(ctx context.Context, id uint) error
	findAllFn  func(ctx context.Context, userID uint, qp dtos.QueryParams) (dtos.TodoItemsPaginatedResponseDto, error)
}

func (r *MockTodoRepository) Create(ctx context.Context, todoItem *todos.TodoItem) error {
	return r.createFn(ctx, todoItem)
}

func (r *MockTodoRepository) FindByID(ctx context.Context, id uint) (todos.TodoItem, error) {
	return r.findByIdFn(ctx, id)
}

func (r *MockTodoRepository) Update(ctx context.Context, id uint, todoItem todos.TodoItem) error {
	return r.updateFn(ctx, id, todoItem)
}

func (r *MockTodoRepository) Delete(ctx context.Context, id uint) error {
	return r.deleteFn(ctx, id)
}

func (r *MockTodoRepository) FindAllForUser(ctx context.Context, userID uint, qp dtos.QueryParams) (dtos.TodoItemsPaginatedResponseDto, error) {
	return r.findAllFn(ctx, userID, qp)
}
