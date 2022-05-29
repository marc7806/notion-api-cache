package filtertreeparser

import (
	"github.com/marc7806/notion-cache/notion"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoDbParser struct{}

func (m MongoDbParser) parse(filterTree *notion.FilterTreeNode) *bson.M {
	if filterTree.CompoundType == nil {
		return &bson.M{
			"properties." + filterTree.Operation.Property + ".value": mapOperationToMongoDbRepresentation(filterTree.Operation),
		}
	}

	var mappedChildren []interface{}
	for _, child := range filterTree.Children {
		mappedChildren = append(mappedChildren, m.parse(&child))
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
	} else if operation.Condition == notion.Contains {
		return bson.M{
			"$regex": operation.Value,
		}
	} else {
		return bson.M{
			"$regex": "",
		}
	}
}
