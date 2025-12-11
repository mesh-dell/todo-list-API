package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mesh-dell/todo-list-API/internal/auth"
	"github.com/mesh-dell/todo-list-API/internal/auth/dtos"
	"github.com/mesh-dell/todo-list-API/internal/auth/service"
	custom "github.com/mesh-dell/todo-list-API/internal/errors"
)

type AuthHandler struct {
	authSvc       *service.AuthService
	refreshSvc    *service.TokenService
	accessSecret  string
	refreshSecret string
	accessExp     int
	refreshExp    int
}

func NewAuthHandler(a *service.AuthService, r *service.TokenService, accessSecret, refreshSecret string, aExp, rExp int) *AuthHandler {
	return &AuthHandler{
		authSvc:       a,
		refreshSvc:    r,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExp:     aExp,
		refreshExp:    rExp,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dtos.RegisterDto
	if err := c.ShouldBind(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.authSvc.Register(req, c.Request.Context())
	if err != nil {
		if errors.Is(err, custom.ErrEmailExists) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// generate access token
	accessToken, refreshToken, refExp := auth.GenerateTokenPair(
		id,
		h.accessSecret,
		h.refreshSecret,
		h.refreshExp,
		h.accessExp,
	)
	// gen ref token
	// save ref token to db
	_ = h.refreshSvc.SaveRefreshToken(id, refreshToken, refExp, c.Request.Context())

	c.IndentedJSON(http.StatusCreated, dtos.AuthResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dtos.LoginDto
	if err := c.ShouldBind(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authSvc.Login(req, c.Request.Context())
	if err != nil {
		if errors.Is(err, custom.ErrInvalidCredentials) {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
	}

	// generate tokens
	accessToken, refreshToken, refExp := auth.GenerateTokenPair(
		user.ID,
		h.accessSecret,
		h.refreshSecret,
		h.refreshExp,
		h.accessExp,
	)
	_ = h.refreshSvc.SaveRefreshToken(user.ID, refreshToken, refExp, c.Request.Context())
	c.IndentedJSON(http.StatusOK, dtos.AuthResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	token := c.GetHeader("Refresh-Token")
	if token == "" {
		c.IndentedJSON(401, gin.H{"error": "missing refresh token"})
		return
	}
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(h.refreshSecret), nil
	})

	if err != nil || !parsed.Valid {
		c.IndentedJSON(401, gin.H{"error": "invalid refresh token"})
		return
	}
	claims := parsed.Claims.(jwt.MapClaims)
	if claims["type"] != "refresh" {
		c.JSON(401, gin.H{"error": "not a refresh token"})
		return
	}
	stored, ok := h.refreshSvc.ValidateRefreshToken(token, c.Request.Context())
	if !ok {
		c.IndentedJSON(401, gin.H{"error": "refresh expired or revoked"})
		return
	}
	accessToken, newRefreshToken, refExp := auth.GenerateTokenPair(
		stored.UserID,
		h.accessSecret,
		h.refreshSecret,
		h.refreshExp,
		h.accessExp,
	)
	_ = h.refreshSvc.RotateRefreshToken(token, newRefreshToken, stored.UserID, refExp, c.Request.Context())
	c.IndentedJSON(200, dtos.AuthResponseDto{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	})
}
