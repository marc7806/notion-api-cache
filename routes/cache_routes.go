package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/cache"
	"github.com/marc7806/notion-cache/config"
	"github.com/marc7806/notion-cache/database"
)

type CacheStatusInformationResponse struct {
	LastUpdated               string `json:"last_updated"`
	NumberOfDatabaseDocuments int32  `json:"num_database_documents"`
}

var (
	refreshChannel = make(chan bool)
)

func AddCacheRoutes(router *gin.RouterGroup) {
	cacheEndpoints := router.Group("/cache")
	{
		cacheEndpoints.POST("/refresh", refreshCacheEndpoint)
		cacheEndpoints.GET("/status", cacheStatusInformationEndpoint)
	}
}

func refreshCacheEndpoint(c *gin.Context) {
	go refreshNotionCache()
	c.JSON(http.StatusOK, gin.H{"status": "Successfully triggered cache refresh"})
}

func cacheStatusInformationEndpoint(c *gin.Context) {
	c.String(200, "Cache not initialize")
}

func refreshNotionCache() {
	refreshNotAlreadyTriggered := <-refreshChannel
	if refreshNotAlreadyTriggered {
		return
	}

	refreshChannel <- true
	client := database.InitClient()

	// Check the connection
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	err = cache.CacheNotionDatabases(client, config.NotionDatabases)
	if err != nil {
		log.Fatal(err)
	}

	// Close the connection once no longer needed
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	refreshChannel <- false
}
