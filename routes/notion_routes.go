package routes

import "github.com/gin-gonic/gin"

func AddNotionRoutes(router *gin.RouterGroup) {
	// todo: add query endpoint
	// todo: create struct for post JSON body data
	// todo: add notion query parser for translation into mongodb query filter

	notionEndpoints := router.Group("/api/notion")
	{
		notionEndpoints.POST("/query")
	}
}
