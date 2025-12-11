package service

import (
	"context"

	"github.com/mesh-dell/todo-list-API/internal/auth"
	"github.com/mesh-dell/todo-list-API/internal/auth/dtos"
	"github.com/mesh-dell/todo-list-API/internal/auth/repository"
	custom "github.com/mesh-dell/todo-list-API/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.IAuthRepository
}

func NewAuthService(r repository.IAuthRepository) *AuthService {
	return &AuthService{
		repo: r,
	}
}

func (svc *AuthService) Login(req dtos.LoginDto, context context.Context) (*auth.User, error) {
	user, _ := svc.repo.GetUserByEmail(context, req.Email)
	if user == nil {
		return nil, custom.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, custom.ErrInvalidCredentials
	}
	return user, nil
}

func (svc *AuthService) Register(req dtos.RegisterDto, context context.Context) (uint, error) {
	if exists, _ := svc.repo.GetUserByEmail(context, req.Email); exists != nil {
		return 0, custom.ErrEmailExists
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &auth.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(bytes),
	}
	err := svc.repo.CreateUser(context, user)
	if err != nil {
		return 0, err
	}
	return user.ID, err
}
