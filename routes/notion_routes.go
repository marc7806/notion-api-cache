package routes

import "github.com/gin-gonic/gin"

func AddNotionRoutes(router *gin.RouterGroup) {
	notionEndpoints := router.Group("/api/notion")
	{
		notionEndpoints.POST("/query")
	}
}
