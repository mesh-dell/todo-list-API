package service

import (
	"context"
	"time"

	"github.com/mesh-dell/todo-list-API/internal/auth"
	"github.com/mesh-dell/todo-list-API/internal/auth/repository"
)

type TokenService struct {
	repo repository.RefreshTokenRepository
}

func NewTokenService(r repository.RefreshTokenRepository) *TokenService {
	return &TokenService{
		repo: r,
	}
}

func (s *TokenService) SaveRefreshToken(userID uint, token string, exp time.Time, ctx context.Context) error {
	rt := &auth.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: exp,
	}
	return s.repo.Save(ctx, rt)
}

func (s *TokenService) ValidateRefreshToken(token string, ctx context.Context) (*auth.RefreshToken, bool) {
	stored, _ := s.repo.Find(ctx, token)
	if stored == nil {
		return nil, false
	}

	// token expired in db
	if time.Now().After(stored.ExpiresAt) {
		s.repo.Delete(ctx, token)
		return nil, false
	}
	return stored, true
}

func (s *TokenService) RotateRefreshToken(oldToken, newToken string, userID uint, newExp time.Time, ctx context.Context) error {
	_ = s.repo.Delete(ctx, oldToken)
	return s.SaveRefreshToken(userID, newToken, newExp, ctx)
}
