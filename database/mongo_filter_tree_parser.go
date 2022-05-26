package database

import (
	"github.com/marc7806/notion-cache/notion"
	"go.mongodb.org/mongo-driver/bson"
)

func ParseToMongoDbQuery(filterTree *notion.FilterTreeNode) *bson.M {
	if filterTree.CompoundType == nil {
		return &bson.M{
			"properties." + filterTree.Operation.Property + ".value": mapOperationToMongoDbRepresentation(filterTree.Operation),
		}
	}

	var mappedChildren []interface{}
	for _, child := range filterTree.Children {
		mappedChildren = append(mappedChildren, ParseToMongoDbQuery(&child))
	}

	return &bson.M{
		string("$" + *filterTree.CompoundType): mappedChildren,
	}
}

func mapOperationToMongoDbRepresentation(operation *notion.FilterOperation) interface{} {
	if operation.Condition == notion.Equals {
		return bson.M{
			"$eq": operation.Value,
		}
	} else {
		return bson.M{}
	}
}
