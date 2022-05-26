package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/database"
	"github.com/marc7806/notion-cache/notion"
)

type QueryRequestBody struct {
	DatabaseId string                 `json:"database_id"`
	Filter     map[string]interface{} `json:"filter"`
	Sorts      []QuerySort            `json:"sorts"`
	Start      int64                  `json:"start"`
	Size       int64                  `json:"page_size"`
}

type QuerySort struct {
	Property  string
	Direction string
}

func AddNotionRoutes(router *gin.RouterGroup) {
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
	filterTree := notion.CreateFilterTree(requestBody.Filter)
	results := database.QueryData(requestBody.DatabaseId, database.ParseToMongoDbQuery(filterTree), requestBody.Sorts, requestBody.Size, requestBody.Start)
	c.JSON(http.StatusOK, results)
}
