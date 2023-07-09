package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/database"
	filtertreeparser "github.com/marc7806/notion-cache/database/filtertree-parser"
	"github.com/marc7806/notion-cache/notion"
	"github.com/marc7806/notion-cache/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddNotionRoutes(router *gin.RouterGroup) {
	notionEndpoints := router.Group("/databases")
	{
		notionEndpoints.POST("/:databaseId/query", queryEndpoint)
	}
}

func queryEndpoint(c *gin.Context) {
	databaseId := c.Param("databaseId")
	requestBody := types.QueryRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var findQuery *primitive.M
	if requestBody.Filter == nil {
		findQuery = &bson.M{}
	} else {
		filterTree := notion.CreateFilterTree(requestBody.Filter)
		findQuery = filtertreeparser.ParseFilterTree(filterTree)
	}

	results, nextCursor, hasMore := database.QueryData(database.DataStore, databaseId, findQuery, requestBody.Sorts, requestBody.Size, requestBody.StartCursor)

	c.JSON(http.StatusOK, types.QueryResponseBody{
		Object:     "list",
		Results:    results,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	})
}
