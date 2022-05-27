package filtertreeparser

import (
	"github.com/marc7806/notion-cache/notion"
	"go.mongodb.org/mongo-driver/bson"
)

type FilterTreeParser interface {
	parse(*notion.FilterTreeNode) *bson.M
}

func ParseFilterTree(filterTreeNode *notion.FilterTreeNode) *bson.M {
	dbParser := MongoDbParser{}
	return dbParser.parse(filterTreeNode)
}
