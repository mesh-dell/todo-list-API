package repository

import (
	"context"

	"github.com/mesh-dell/todo-list-API/internal/todos"
	"github.com/mesh-dell/todo-list-API/internal/utils"
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
func (r *TodoRepository) FindAllForUser(ctx context.Context, userID uint, pagination *utils.Pagination) (*utils.Pagination, error) {
	items, err := gorm.G[todos.TodoItem](r.db).Where("user_id = ?", userID).Find(ctx)
	if err != nil {
		return nil, err
	}
	p, err := utils.Paginate(r.db, ctx, pagination, &items)
	if err != nil {
		return nil, err
	}
	return p, nil
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
