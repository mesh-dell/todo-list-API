package main

import (
	"log"

	"github.com/mesh-dell/todo-list-API/config"
	"github.com/mesh-dell/todo-list-API/internal/api"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	api.InitServer(config)
}
