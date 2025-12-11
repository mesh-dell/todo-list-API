package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mesh-dell/todo-list-API/config"
	"github.com/mesh-dell/todo-list-API/internal/auth"
	"github.com/mesh-dell/todo-list-API/internal/auth/handler"
	"github.com/mesh-dell/todo-list-API/internal/auth/repository"
	"github.com/mesh-dell/todo-list-API/internal/auth/service"
	"github.com/mesh-dell/todo-list-API/internal/database"
	"github.com/mesh-dell/todo-list-API/internal/todos"
)

func InitServer(config config.Config) {
	gormDb := database.NewGormDb(config)
	gormDb.DbClient.AutoMigrate(&auth.User{}, &todos.TodoItem{}, &auth.RefreshToken{})

	authRepository := repository.NewAuthRepository(gormDb.DbClient)
	authService := service.NewAuthService(authRepository)
	tokenRepository := repository.NewRefreshTokenRepository(gormDb.DbClient)
	tokenService := service.NewTokenService(*tokenRepository)
	authHandler := handler.NewAuthHandler(
		authService,
		tokenService,
		os.Getenv("ACCESS_SECRET"),
		os.Getenv("REFRESH_SECRET"),
		config.JWTExpiry,
		config.RefreshExpiry,
	)
	router := gin.Default()
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/refresh", authHandler.Refresh)
	router.Run()
}
