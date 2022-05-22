package notion

import (
	"fmt"
	"time"
)

type Page struct {
	ID             string    `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTime    time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	CreatedBy      string
	LastEditedBy   string
	Properties     []interface{}
}

type SinglePageProperty struct {
	Name  string
	Type  string
	Value string
}

type MultiPageProperty struct {
	Name  string
	Type  string
	Value []string
}

func ParsePage(notionResponseObject *NotionDatabaseObject) *Page {
	// page := Page{
	// 	ID:             notionResponseObject.ID,
	// 	CreatedTime:    notionResponseObject.CreatedTime,
	// 	LastEditedTime: notionResponseObject.LastEditedTime,
	// 	CreatedBy:      notionResponseObject.CreatedBy.ID,
	// 	LastEditedBy:   notionResponseObject.LastEditedBy.ID,
	// }

	for key, prop := range notionResponseObject.Properties {
		// fmt.Println(key, prop)
		fmt.Println(key)
		fmt.Println(ResolvePropertyType(prop))

	}

	return nil
}
