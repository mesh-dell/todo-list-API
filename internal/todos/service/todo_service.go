package service

import (
	"context"
	"errors"

	custom "github.com/mesh-dell/todo-list-API/internal/errors"
	"github.com/mesh-dell/todo-list-API/internal/todos"
	"github.com/mesh-dell/todo-list-API/internal/todos/dtos"
	"github.com/mesh-dell/todo-list-API/internal/todos/repository"
	"gorm.io/gorm"
)

type TodoService struct {
	repo repository.ITodoRepository
}

func NewTodoService(r repository.ITodoRepository) *TodoService {
	return &TodoService{
		repo: r,
	}
}

func (svc *TodoService) Create(ctx context.Context, userID uint, req dtos.TodoItemRequestDto) (*todos.TodoItem, error) {
	todoItem := &todos.TodoItem{
		Title:       req.Title,
		Description: req.Description,
		UserId:      userID,
	}
	err := svc.repo.Create(ctx, todoItem)
	return todoItem, err
}

func (svc *TodoService) FindAllForUser(ctx context.Context, userID uint, queryParams dtos.QueryParams) (dtos.TodoItemsPaginatedResponseDto, error) {
	return svc.repo.FindAllForUser(ctx, userID, queryParams)
}

func (svc *TodoService) FindByID(ctx context.Context, id uint) (todos.TodoItem, error) {
	return svc.repo.FindByID(ctx, id)
}

func (svc *TodoService) Delete(ctx context.Context, id uint, userID uint) error {
	item, err := svc.FindByID(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return custom.ErrItemNotFound
		}
	}

	if item.UserId != userID {
		return custom.ErrCannotDeleteItem
	}

	return svc.repo.Delete(ctx, id)
}

func (svc *TodoService) Update(ctx context.Context, id uint, userID uint, req dtos.TodoItemRequestDto) (todos.TodoItem, error) {
	item, err := svc.FindByID(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return todos.TodoItem{}, custom.ErrItemNotFound
		}
	}

	if item.UserId != userID {
		return todos.TodoItem{}, custom.ErrCannotUpdateItem
	}

	item.Title = req.Title
	item.Description = req.Description

	err = svc.repo.Update(ctx, id, item)
	if err != nil {
		return todos.TodoItem{}, err
	}
	return item, nil
}
