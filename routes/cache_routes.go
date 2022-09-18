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
	Status                    string `json:"status"`
}

func AddCacheRoutes(router *gin.RouterGroup) {
	cacheEndpoints := router.Group("/cache")
	{
		cacheEndpoints.POST("/refresh", refreshCacheEndpoint)
		cacheEndpoints.GET("/status", cacheStatusInformationEndpoint)
		cacheEndpoints.POST("/clear", clearCacheEndpoint)
	}
}

func refreshCacheEndpoint(c *gin.Context) {
	isRefreshing := cache.HandleCacheRefresh()
	if isRefreshing {
		c.JSON(http.StatusAccepted, gin.H{"status": "Cache is currently refreshing"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Successfully triggered new cache refresh"})
}

func clearCacheEndpoint(c *gin.Context) {
	isRefreshing := cache.HandleCacheClear()
	if isRefreshing {
		c.JSON(http.StatusAccepted, gin.H{"status": "Cache is currently refreshing"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Successfully triggered cache clearing"})
}

func cacheStatusInformationEndpoint(c *gin.Context) {
	cacheInfo := CacheStatusInformationResponse{
		LastUpdated:               cache.LastUpdated.Format(time.RFC3339),
		NumberOfDatabaseDocuments: cache.NumUpdatedDocs,
		Status:                    string(cache.Status),
	}
	c.JSON(http.StatusAccepted, cacheInfo)
}
