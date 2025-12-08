package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mesh-dell/todo-list-API/config"
	"github.com/mesh-dell/todo-list-API/internal/auth"
	"github.com/mesh-dell/todo-list-API/internal/database"
	"github.com/mesh-dell/todo-list-API/internal/todos"
)

func InitServer(config config.Config) {
	gormDb := database.NewGormDb(config)
	gormDb.DbClient.AutoMigrate(&auth.User{}, &todos.TodoItem{})

	router := gin.Default()
	router.Run()
}
