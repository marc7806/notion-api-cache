package notion

import (
	"time"
)

type Page struct {
	ID             string    `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTime    time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	CreatedBy      struct {
		Object string `json:"object"`
		ID     string `json:"id"`
	} `json:"created_by"`
	LastEditedBy struct {
		Object string `json:"object"`
		ID     string `json:"id"`
	} `json:"last_edited_by"`
	Cover    interface{} `json:"cover"`
	Icon     interface{} `json:"icon"`
	Archived bool        `json:"archived"`
	URL      string      `json:"url"`
	Parent   struct {
		Type       string `json:"type"`
		DatabaseID string `json:"database_id"`
	} `json:"parent"`
	Properties []PageProperty
}

type PageProperty struct {
	Name  string
	Type  string
	Value string
}

func ParsePage(notionResponseObject *NotionDatabaseObject) *Page {
	page := Page{
		ID:             notionResponseObject.ID,
		CreatedTime:    notionResponseObject.CreatedTime,
		LastEditedTime: notionResponseObject.LastEditedTime,
		CreatedBy:      notionResponseObject.CreatedBy.ID,
		LastEditedBy:   notionResponseObject.LastEditedBy.ID,
	}
	var pageProps []PageProperty
	for key, prop := range notionResponseObject.Properties {
		propValue, err := ResolvePropertyType(prop)
		if err == nil {
			pageProps = append(pageProps, PageProperty{
				Name:  key,
				Type:  prop["type"].(string),
				Value: propValue,
			})
		}
	}
	page.Properties = pageProps
	return &page
}
