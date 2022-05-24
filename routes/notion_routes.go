package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/notion"
)

type QueryRequestBody struct {
	DatabaseId string                 `json:"database_id"`
	Filter     map[string]interface{} `json:"filter"`
	Sorts      []QuerySort            `json:"sorts"`
	Page       string                 `json:"page"`
	Size       string                 `json:"size"`
}

type QuerySort struct {
	Property  string
	Direction string
}

func AddNotionRoutes(router *gin.RouterGroup) {
	// todo: add query endpoint
	// todo: create struct for post JSON body data
	// todo: add notion query parser for translation into mongodb query filter

	notionEndpoints := router.Group("/notion")
	{
		notionEndpoints.POST("/query", queryEndpoint)
	}
}

func queryEndpoint(c *gin.Context) {
	requestBody := QueryRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	notion.CreateFilterTree(requestBody.Filter)

	c.JSON(http.StatusOK, gin.H{"status": "Successfully triggered new cache refresh"})
}
