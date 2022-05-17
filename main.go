package main

import (
	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/config"
	"github.com/marc7806/notion-cache/routes"
)

var (
	router = gin.Default()
)

func main() {
	config.LoadEnvironment()
	buildRoutes()
	router.Run(":8090")
}

func buildRoutes() {
	api := router.Group("/api")
	routes.AddCacheRoutes(api)
	routes.AddNotionRoutes(api)
}
