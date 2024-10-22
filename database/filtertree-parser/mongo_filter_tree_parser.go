package filtertreeparser

import (
	"github.com/marc7806/notion-cache/notion"
	"github.com/marc7806/notion-cache/utils"
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
	} else if operation.Condition == notion.DoesNotEqual {
		return bson.M{
			"$ne": operation.Value,
		}
	} else if operation.Condition == notion.Contains {
		return bson.M{
			"$regex": operation.Value,
		}
	} else if operation.Condition == notion.StartsWith {
		return bson.M{
			"$regex": "^" + operation.Value,
		}
	} else if operation.Condition == notion.EndsWith {
		return bson.M{
			"$regex": operation.Value + "$",
		}
	} else if operation.Condition == notion.IsNotEmpty {
		boolValue, err := utils.StringToBool(operation.Value)

		if err != nil && boolValue {
			return bson.M{
				"$ne": "",
			}
		} else if err != nil && !boolValue {
			return bson.M{
				"$eq": "",
			}
		}
	} else if operation.Condition == notion.IsEmpty {
		boolValue, err := utils.StringToBool(operation.Value)

		if err != nil && boolValue {
			return bson.M{
				"$eq": "",
			}
		} else if err != nil && !boolValue {
			return bson.M{
				"$ne": "",
			}
		}
	}

	return bson.M{
		"$regex": "",
	}
}
