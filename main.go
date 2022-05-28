package main

import (
	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/routes"
	"github.com/marc7806/notion-cache/scheduler"
)

var (
	router = gin.Default()
)

func main() {
	// todo: add api-token middleware for authentication
	go scheduler.Init()
	buildRoutes()
	router.Run(":8080")
}

func buildRoutes() {
	api := router.Group("/v1")
	routes.AddCacheRoutes(api)
	routes.AddNotionRoutes(api)
}
