package routes

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/cache"
	"github.com/marc7806/notion-cache/config"
	"github.com/marc7806/notion-cache/database"
)

type RefreshState struct {
	mu           sync.Mutex
	isRefreshing bool
}

func (s *RefreshState) setRefreshState(state bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isRefreshing = state
}

type CacheStatusInformationResponse struct {
	LastUpdated               string `json:"last_updated"`
	NumberOfDatabaseDocuments int    `json:"num_database_documents"`
}

var (
	lastUpdated    time.Time
	numUpdatedDocs int
	refreshState   *RefreshState
)

func AddCacheRoutes(router *gin.RouterGroup) {
	cacheEndpoints := router.Group("/cache")
	{
		cacheEndpoints.POST("/refresh", refreshCacheEndpoint)
		cacheEndpoints.GET("/status", cacheStatusInformationEndpoint)
	}
	// initialize refresh state
	refreshState = new(RefreshState)
}

func refreshCacheEndpoint(c *gin.Context) {
	log.Print(refreshState.isRefreshing)
	if refreshState.isRefreshing {
		c.JSON(http.StatusAccepted, gin.H{"status": "Cache is currently refreshing"})
		return
	}
	go refreshNotionCache()
	c.JSON(http.StatusOK, gin.H{"status": "Successfully triggered new cache refresh"})
}

func cacheStatusInformationEndpoint(c *gin.Context) {
	cacheInfo := CacheStatusInformationResponse{
		LastUpdated:               lastUpdated.Format(time.RFC3339),
		NumberOfDatabaseDocuments: numUpdatedDocs,
	}
	c.JSON(http.StatusAccepted, cacheInfo)
}

func refreshNotionCache() {
	refreshState.setRefreshState(true)
	client := database.InitClient()

	// Check the connection
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	updatedDocsLength, err := cache.CacheNotionDatabases(client, config.NotionDatabases)
	if err != nil {
		log.Fatal(err)
	}

	// Close the connection once no longer needed
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	lastUpdated = time.Now()
	numUpdatedDocs = updatedDocsLength
	refreshState.setRefreshState(false)
}
