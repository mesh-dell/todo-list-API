package repository

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/mesh-dell/todo-list-API/internal/todos"
	"github.com/mesh-dell/todo-list-API/internal/todos/dtos"
	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

// Create implements ITodoRepository.
func (r *TodoRepository) Create(ctx context.Context, todoItem *todos.TodoItem) error {
	return gorm.G[todos.TodoItem](r.db).Create(ctx, todoItem)
}

// Delete implements ITodoRepository.
func (r *TodoRepository) Delete(ctx context.Context, id uint) error {
	_, err := gorm.G[todos.TodoItem](r.db).Where("id = ?", id).Delete(ctx)
	return err
}

// FindAllForUser implements ITodoRepository.
func (r *TodoRepository) FindAllForUser(ctx context.Context, userID uint, queryParams dtos.QueryParams) (dtos.TodoItemsPaginatedResponseDto, error) {
	if queryParams.Limit < 1 {
		queryParams.Limit = 10
	}
	if queryParams.Page < 1 {
		queryParams.Page = 1
	}
	offset := (queryParams.Page - 1) * queryParams.Limit

	var count int64
	r.db.Model(&todos.TodoItem{}).Where("user_id = ?", userID).Count(&count)
	totalPages := math.Ceil(float64(count) / float64(queryParams.Limit))

	q := "%" + queryParams.SearchQuery + "%"
	allowedColumn := map[string]bool{
		"created_at":  true,
		"title":       true,
		"description": true,
	}
	if !allowedColumn[strings.ToLower(queryParams.SortBy)] {
		queryParams.SortBy = "created_at"
	}
	if !strings.EqualFold(queryParams.Order, "asc") && !strings.EqualFold(queryParams.Order, "desc") {
		queryParams.Order = "desc"
	}

	orderBy := fmt.Sprintf("%s %s", queryParams.SortBy, queryParams.Order)

	items, err := gorm.G[todos.TodoItem](r.db).
		Offset(offset).
		Limit(queryParams.Limit).Where("user_id = ?", userID).
		Where("title LIKE ? or description LIKE ?", q, q).
		Order(orderBy).
		Find(ctx)

	if err != nil {
		return dtos.TodoItemsPaginatedResponseDto{}, err
	}

	var todoItemsRes []dtos.TodoItemResponseDto
	for _, item := range items {
		itemRes := dtos.TodoItemResponseDto{
			Id:          item.ID,
			Title:       item.Title,
			Description: item.Description,
		}
		todoItemsRes = append(todoItemsRes, itemRes)
	}

	paginatedRes := dtos.TodoItemsPaginatedResponseDto{
		Data:  todoItemsRes,
		Page:  queryParams.Page,
		Limit: queryParams.Limit,
		Total: int(totalPages),
	}
	return paginatedRes, nil
}

// FindByID implements ITodoRepository.
func (r *TodoRepository) FindByID(ctx context.Context, id uint) (todos.TodoItem, error) {
	return gorm.G[todos.TodoItem](r.db).Where("id = ?", id).First(ctx)
}

// Update implements ITodoRepository.
func (r *TodoRepository) Update(ctx context.Context, id uint, todoItem todos.TodoItem) error {
	_, err := gorm.G[todos.TodoItem](r.db).Where("id = ?", id).Updates(ctx, todoItem)
	return err
}

func NewTodoRepository(db *gorm.DB) ITodoRepository {
	return &TodoRepository{
		db: db,
	}
}
