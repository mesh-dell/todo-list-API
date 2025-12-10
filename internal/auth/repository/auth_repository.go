package repository

import (
	"context"

	"github.com/mesh-dell/todo-list-API/internal/auth"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

// CreateUser implements IAuthRepository.
func (r *AuthRepository) CreateUser(context context.Context, user *auth.User) error {
	return gorm.G[auth.User](r.db).Create(context, user)
}

// GetUserByEmail implements IAuthRepository.
func (r *AuthRepository) GetUserByEmail(context context.Context, email string) (*auth.User, error) {
	user, err := gorm.G[auth.User](r.db).Where("email = ?", email).First(context)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &AuthRepository{
		db: db,
	}
}
