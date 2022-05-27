package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marc7806/notion-cache/database"
	filtertreeparser "github.com/marc7806/notion-cache/database/filtertree-parser"
	"github.com/marc7806/notion-cache/notion"
)

type QueryRequestBody struct {
	Filter      map[string]interface{} `json:"filter"`
	Sorts       []QuerySort            `json:"sorts"`
	StartCursor string                 `json:"start_cursor"`
	Size        int64                  `json:"page_size"`
}

type QueryResponseBody struct {
	Object     string         `json:"object"`
	Results    []*notion.Page `json:"results"`
	NextCursor string         `json:"next_cursor"`
	HasMore    bool           `json:"has_more"`
}

type QuerySort struct {
	Property  string
	Direction string
}

func AddNotionRoutes(router *gin.RouterGroup) {
	notionEndpoints := router.Group("/databases")
	{
		notionEndpoints.POST("/:databaseId/query", queryEndpoint)
	}
}

// TODO: add support for cursor based pagination...? to remain constant with the original notion query api
func queryEndpoint(c *gin.Context) {
	databaseId := c.Param("databaseId")
	requestBody := QueryRequestBody{}

	if err := c.Bind(&requestBody); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filterTree := notion.CreateFilterTree(requestBody.Filter)
	results, nextCursor, hasMore := database.QueryData(databaseId, filtertreeparser.ParseFilterTree(filterTree), requestBody.Sorts, requestBody.Size, requestBody.StartCursor)
	c.JSON(http.StatusOK, QueryResponseBody{
		Object:     "list",
		Results:    results,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	})
}
