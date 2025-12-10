package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mesh-dell/todo-list-API/config"
)

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint) (string, error) {
	secret := os.Getenv("ACCESS_SECRET")
	exp := time.Minute * 15
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(exp).Unix(),
		},
	})
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string, config config.Config) (*JWTClaims, error) {
	secret := os.Getenv("ACCESS_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
