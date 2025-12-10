package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mesh-dell/todo-list-API/internal/auth/dtos"
	"github.com/mesh-dell/todo-list-API/internal/auth/service"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: s,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dtos.RegisterDto
	if err := c.ShouldBind(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.svc.Register(req, c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, dtos.AuthResponseDto{
		Token: token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dtos.LoginDto
	if err := c.ShouldBind(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.svc.Login(req, c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, dtos.AuthResponseDto{
		Token: token,
	})
}
