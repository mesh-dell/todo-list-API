package repository

import (
	"context"
	"errors"

	"github.com/mesh-dell/todo-list-API/internal/auth"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		DB: db,
	}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, token *auth.RefreshToken) error {
	return gorm.G[auth.RefreshToken](r.DB).Create(ctx, token)
}

func (r *RefreshTokenRepository) Find(ctx context.Context, tokenString string) (*auth.RefreshToken, error) {
	token, err := gorm.G[auth.RefreshToken](r.DB).Where("token = ?", tokenString).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &token, err
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, tokenString string) error {
	_, err := gorm.G[auth.RefreshToken](r.DB).Where("token = ?", tokenString).Delete(ctx)
	return err
}

func (r RefreshTokenRepository) DeleteAllTokensForUser(ctx context.Context, userID uint) error {
	_, err := gorm.G[auth.RefreshToken](r.DB).Where("user_id = ?", userID).Delete(ctx)
	return err
}
