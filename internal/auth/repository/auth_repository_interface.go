package repository

import (
	"context"

	"github.com/mesh-dell/todo-list-API/internal/auth"
)

type IAuthRepository interface {
	CreateUser(context context.Context, user *auth.User) error
	GetUserByEmail(context context.Context, email string) (*auth.User, error)
}
