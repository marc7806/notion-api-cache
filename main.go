package main

import (
	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/cache"
	"github.com/marc7806/notion-cache/database"
	"github.com/marc7806/notion-cache/routes"
	"github.com/marc7806/notion-cache/scheduler"
)

var (
	router = gin.Default()
)

func main() {
	// todo: add api-token middleware for authentication

	// initialize database
	datastore, err := database.InitMongoDataStore()
	if err != nil {
		panic(err)
	}
	// disconnect database client once application shuts down
	defer database.DisconnectDb(datastore)

	// initialize caching metadata store
	cache.InitCache()
	// init caching scheduler to be run in defined interval
	go scheduler.InitScheduler()

	buildRoutes()
	router.Run(":8080")
}

func buildRoutes() {
	api := router.Group("/v1")
	routes.AddCacheRoutes(api)
	routes.AddNotionRoutes(api)
}
