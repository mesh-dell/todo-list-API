package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateTokenPair(userID uint, accessSecret, refreshSecret string, refreshExp, accessExp int) (string, string, time.Time) {
	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"type": "access",
		"exp":  time.Now().Add(time.Duration(accessExp) * time.Second).Unix(),
	}
	access, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(accessSecret))

	// refresh token
	refreshExpiresAt := time.Now().Add(time.Duration(refreshExp) * time.Second)
	refreshClaims := jwt.MapClaims{
		"sub":  userID,
		"type": "refresh",
		"exp":  refreshExpiresAt.Unix(),
	}
	refresh, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(refreshSecret))
	return access, refresh, refreshExpiresAt
}

func ValidateJWT(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
