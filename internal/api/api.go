package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mesh-dell/todo-list-API/config"
	"github.com/mesh-dell/todo-list-API/internal/auth"
	"github.com/mesh-dell/todo-list-API/internal/auth/handler"
	"github.com/mesh-dell/todo-list-API/internal/auth/middleware"
	"github.com/mesh-dell/todo-list-API/internal/auth/repository"
	"github.com/mesh-dell/todo-list-API/internal/auth/service"
	"github.com/mesh-dell/todo-list-API/internal/database"
	"github.com/mesh-dell/todo-list-API/internal/todos"
	todoHandler "github.com/mesh-dell/todo-list-API/internal/todos/handler"
	todoRepository "github.com/mesh-dell/todo-list-API/internal/todos/repository"
	todoService "github.com/mesh-dell/todo-list-API/internal/todos/service"
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

	todoRepository := todoRepository.NewTodoRepository(gormDb.DbClient)
	todoService := todoService.NewTodoService(todoRepository)
	todoHandler := todoHandler.NewTodoHandler(*todoService)

	router := gin.Default()
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/refresh", authHandler.Refresh)

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(os.Getenv("ACCESS_SECRET")))
	{
		protected.POST("", todoHandler.Create)
		protected.GET("/:id", todoHandler.FindByID)
		protected.GET("", todoHandler.FindAllForUser)
		protected.PUT("/:id", todoHandler.Update)
		protected.DELETE("/:id", todoHandler.Delete)
	}
	router.Run()
}
