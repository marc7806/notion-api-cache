package sortparser

import (
	"github.com/marc7806/notion-cache/types"
	"github.com/marc7806/notion-cache/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func ParseSortOptions(sort []types.QuerySort) *bson.M {
	result := bson.M{}
	if len(sort) == 0 {
		// default sort by id
		return &bson.M{"_id": 1}
	}

	for _, sortOption := range sort {
		result["properties."+sortOption.Property+".value"] = utils.BinarySortDirection(sortOption.Direction)
	}

	return &result
}
