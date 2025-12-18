package tests

import (
	"context"
	"errors"
	"testing"

	custom "github.com/mesh-dell/todo-list-API/internal/errors"
	"github.com/mesh-dell/todo-list-API/internal/todos"
	"github.com/mesh-dell/todo-list-API/internal/todos/dtos"
	"github.com/mesh-dell/todo-list-API/internal/todos/service"
	"gorm.io/gorm"
)

func TestTodoService_Create_Success(t *testing.T) {
	repo := &MockTodoRepository{
		createFn: func(ctx context.Context, todoItem *todos.TodoItem) error {
			return nil
		},
	}
	svc := service.NewTodoService(repo)

	item, err := svc.Create(context.Background(), 1, dtos.TodoItemRequestDto{
		Title:       "test",
		Description: "test",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item.UserId != 1 {
		t.Fatalf("userID not set correctly")
	}
}

func TestTodoService_FindByID_NotFound(t *testing.T) {
	repo := &MockTodoRepository{
		findByIdFn: func(ctx context.Context, id uint) (todos.TodoItem, error) {
			return todos.TodoItem{}, gorm.ErrRecordNotFound
		},
	}

	svc := service.NewTodoService(repo)
	_, err := svc.FindByID(context.Background(), 1, 1)

	if !errors.Is(err, custom.ErrItemNotFound) {
		t.Fatalf("expected ErrItemNotFound, got %v", err)
	}
}

func TestTodoService_FindByID_WrongUser(t *testing.T) {
	repo := &MockTodoRepository{
		findByIdFn: func(ctx context.Context, id uint) (todos.TodoItem, error) {
			return todos.TodoItem{
				UserId: 1,
			}, nil
		},
	}

	svc := service.NewTodoService(repo)
	_, err := svc.FindByID(context.Background(), 1, 2)

	if !errors.Is(err, custom.ErrCannotGetItem) {
		t.Fatalf("expected ErrCannotGetItem, got %v", err)
	}
}

func TestTodoService_FindByID_Success(t *testing.T) {
	repo := &MockTodoRepository{
		findByIdFn: func(ctx context.Context, id uint) (todos.TodoItem, error) {
			return todos.TodoItem{
				UserId: 1,
				Title:  "todo",
			}, nil
		},
	}

	svc := service.NewTodoService(repo)
	item, err := svc.FindByID(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if item.Title != "todo" {
		t.Fatal("unexpected item returned")
	}
}

func TestTodoService_Delete_WrongUser(t *testing.T) {
	repo := &MockTodoRepository{
		findByIdFn: func(ctx context.Context, id uint) (todos.TodoItem, error) {
			return todos.TodoItem{
				UserId: 1,
			}, nil
		},
	}

	svc := service.NewTodoService(repo)
	err := svc.Delete(context.Background(), 1, 5)

	if !errors.Is(err, custom.ErrCannotDeleteItem) {
		t.Fatalf("expected ErrCannotDeleteItem, got %v", err)
	}
}

func TestTodoService_Delete_Success(t *testing.T) {
	deleted := false
	repo := &MockTodoRepository{
		findByIdFn: func(ctx context.Context, id uint) (todos.TodoItem, error) {
			return todos.TodoItem{
				UserId: 1,
			}, nil
		},
		deleteFn: func(ctx context.Context, id uint) error {
			deleted = true
			return nil
		},
	}

	svc := service.NewTodoService(repo)
	err := svc.Delete(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if !deleted {
		t.Fatalf("expected delete to be called")
	}
}

func TestTodoService_Update_WrongUser(t *testing.T) {
	repo := &MockTodoRepository{
		findByIdFn: func(ctx context.Context, id uint) (todos.TodoItem, error) {
			return todos.TodoItem{
				UserId: 1,
			}, nil
		},
	}

	svc := service.NewTodoService(repo)
	_, err := svc.Update(context.Background(), 1, 5, dtos.TodoItemRequestDto{})

	if !errors.Is(err, custom.ErrCannotUpdateItem) {
		t.Fatalf("expected ErrCannotUpdateItem, got %v", err)
	}
}

func TestTodoService_Update_Success(t *testing.T) {
	repo := &MockTodoRepository{
		findByIdFn: func(ctx context.Context, id uint) (todos.TodoItem, error) {
			return todos.TodoItem{
				UserId: 1,
				Title:  "todo",
			}, nil
		},
		updateFn: func(ctx context.Context, id uint, todoItem todos.TodoItem) error {
			return nil
		},
	}

	svc := service.NewTodoService(repo)
	item, err := svc.Update(context.Background(), 1, 1, dtos.TodoItemRequestDto{Title: "changed"})

	if err != nil {
		t.Fatalf("unexpected error: %v", item)
	}

	if item.Title != "changed" {
		t.Fatalf("item not updated")
	}
}
