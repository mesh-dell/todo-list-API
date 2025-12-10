package service

import (
	"context"
	"fmt"

	"github.com/mesh-dell/todo-list-API/internal/auth"
	"github.com/mesh-dell/todo-list-API/internal/auth/dtos"
	"github.com/mesh-dell/todo-list-API/internal/auth/repository"
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

func (svc *AuthService) Login(req dtos.LoginDto, context context.Context) (string, error) {
	user, err := svc.repo.GetUserByEmail(context, req.Email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	return GenerateJWT(user.ID)
}

func (svc *AuthService) Register(req dtos.RegisterDto, context context.Context) (string, error) {
	if existing, _ := svc.repo.GetUserByEmail(context, req.Email); existing != nil {
		return "", fmt.Errorf("email already exists")
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &auth.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(bytes),
	}
	err := svc.repo.CreateUser(context, user)
	if err != nil {
		return "", err
	}
	token, err := GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
