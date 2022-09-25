package types

import "github.com/marc7806/notion-cache/notion"

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
	Direction QuerySortDirection
}

type QuerySortDirection string

const (
	Asc  QuerySortDirection = "ascending"
	Desc                    = "descending"
)
