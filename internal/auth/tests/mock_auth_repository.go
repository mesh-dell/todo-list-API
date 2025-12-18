package tests

import (
	"context"

	"github.com/mesh-dell/todo-list-API/internal/auth"
)

type MockAuthRepository struct {
	getUserByEmailFn func(context context.Context, email string) (*auth.User, error)
	createUser       func(context context.Context, user *auth.User) error
}

func (r *MockAuthRepository) GetUserByEmail(context context.Context, email string) (*auth.User, error) {
	return r.getUserByEmailFn(context, email)
}

func (r *MockAuthRepository) CreateUser(context context.Context, user *auth.User) error {
	return r.createUser(context, user)
}
