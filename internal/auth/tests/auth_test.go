package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/mesh-dell/todo-list-API/internal/auth"
	"github.com/mesh-dell/todo-list-API/internal/auth/dtos"
	"github.com/mesh-dell/todo-list-API/internal/auth/service"
	custom "github.com/mesh-dell/todo-list-API/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := &MockAuthRepository{
		getUserByEmailFn: func(context context.Context, email string) (*auth.User, error) {
			return nil, nil
		},
	}
	svc := service.NewAuthService(repo)
	_, err := svc.Login(dtos.LoginDto{Email: "test@email.com", Password: "password"}, context.Background())

	if !errors.Is(err, custom.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials got %v", err)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("correctPassword"), bcrypt.DefaultCost)
	repo := &MockAuthRepository{
		getUserByEmailFn: func(context context.Context, email string) (*auth.User, error) {
			return &auth.User{
				Email:        email,
				PasswordHash: string(passwordHash),
			}, nil
		},
	}
	svc := service.NewAuthService(repo)

	_, err := svc.Login(dtos.LoginDto{
		Email:    "test@email.com",
		Password: "wrongPassword",
	}, context.Background())

	if !errors.Is(err, custom.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials got %v", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	repo := &MockAuthRepository{
		getUserByEmailFn: func(context context.Context, email string) (*auth.User, error) {
			return &auth.User{
				Email:        email,
				PasswordHash: string(hash),
			}, nil
		},
	}

	svc := service.NewAuthService(repo)
	user, err := svc.Login(dtos.LoginDto{Email: "test@email.com", Password: "password"}, context.Background())

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if user.Email != "test@email.com" {
		t.Fatalf("expected email test@email.com, got %v", user.Email)
	}
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	repo := &MockAuthRepository{
		getUserByEmailFn: func(context context.Context, email string) (*auth.User, error) {
			return &auth.User{
				Email: email,
			}, nil
		},
	}
	svc := service.NewAuthService(repo)
	_, err := svc.Register(dtos.RegisterDto{
		Name:     "John",
		Email:    "test@email.com",
		Password: "password",
	}, context.Background())

	if !errors.Is(err, custom.ErrEmailExists) {
		t.Fatalf("expected ErrEmailExists got %v", err)
	}
}

func TestAuthService_Register_CreateError(t *testing.T) {
	repo := &MockAuthRepository{
		createUser: func(context context.Context, user *auth.User) error {
			return fmt.Errorf("db error")
		},
		getUserByEmailFn: func(context context.Context, email string) (*auth.User, error) {
			return nil, nil
		},
	}

	svc := service.NewAuthService(repo)
	_, err := svc.Register(dtos.RegisterDto{
		Name:     "John",
		Email:    "test@email.com",
		Password: "password",
	}, context.Background())

	if err == nil {
		t.Fatalf("expected error, returned nil")
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	repo := &MockAuthRepository{
		getUserByEmailFn: func(context context.Context, email string) (*auth.User, error) {
			return nil, nil
		},
		createUser: func(context context.Context, user *auth.User) error {
			user.ID = 1
			return nil
		},
	}

	svc := service.NewAuthService(repo)
	id, err := svc.Register(dtos.RegisterDto{
		Name:     "john",
		Email:    "test@email.com",
		Password: "password",
	}, context.Background())

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if id != 1 {
		t.Fatalf("expected userID 1, got %d", id)
	}
}
