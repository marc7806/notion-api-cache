package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/cache"
)

type CacheStatusInformationResponse struct {
	LastUpdated               string `json:"last_updated"`
	NumberOfDatabaseDocuments int    `json:"num_database_documents"`
}

var (
	lastUpdated    time.Time
	numUpdatedDocs int
)

func AddCacheRoutes(router *gin.RouterGroup) {
	cacheEndpoints := router.Group("/cache")
	{
		cacheEndpoints.POST("/refresh", refreshCacheEndpoint)
		cacheEndpoints.GET("/status", cacheStatusInformationEndpoint)
	}
}

func refreshCacheEndpoint(c *gin.Context) {
	isRefreshing, lastUpdatedTime, numDocs := cache.HandleCacheRefresh()
	lastUpdated = lastUpdatedTime
	numUpdatedDocs = numDocs
	if isRefreshing {
		c.JSON(http.StatusAccepted, gin.H{"status": "Cache is currently refreshing"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Successfully triggered new cache refresh"})
}

func cacheStatusInformationEndpoint(c *gin.Context) {
	cacheInfo := CacheStatusInformationResponse{
		LastUpdated:               lastUpdated.Format(time.RFC3339),
		NumberOfDatabaseDocuments: numUpdatedDocs,
	}
	c.JSON(http.StatusAccepted, cacheInfo)
}
