package notion

type FilterTreeNode struct {
	CompoundType *CompoundFilterType
	Operation    *FilterOperation
	Children     []FilterTreeNode
}

type FilterOperation struct {
	Property  string
	Condition FilterCondition
	Value     string
}

type CompoundFilterType string

const (
	Leaf CompoundFilterType = "leaf"
	And                     = "and"
	Or                      = "or"
)

type FilterCondition int64

const (
	Equals FilterCondition = iota
	IsNotEmpty
)

func CreateFilterTree(filterQuery map[string]interface{}) *FilterTreeNode {
	// recursive parse query body from rest endpoint using dfs and return filter tree node
	var filterTreeNode *FilterTreeNode
	return traverseFilterQuery(filterQuery, filterTreeNode)
}

func traverseFilterQuery(filterQuery map[string]interface{}, filterTreeNode *FilterTreeNode) *FilterTreeNode {
	// 1. case: single property filter - Basecase
	if isPropertyField(filterQuery) {
		condition, value := resolveFilterCondition(filterQuery)
		filterOp := FilterOperation{
			Property:  filterQuery["property"].(string),
			Condition: condition,
			Value:     value,
		}
		return &FilterTreeNode{
			Operation: &filterOp,
		}
	}

	// 2. case: nested multi filter
	filterCompoundType := resolveCompoundFilterType(filterQuery)
	if filterTreeNode == nil {
		filterTreeNode = &FilterTreeNode{
			CompoundType: &filterCompoundType,
		}
	}
	for _, el := range filterQuery[string(filterCompoundType)].([]interface{}) {
		filterTreeNode.Children = append(filterTreeNode.Children, *traverseFilterQuery(el.(map[string]interface{}), filterTreeNode))
	}
	return filterTreeNode
}

func resolveFilterCondition(propertyFilterObject map[string]interface{}) (condition FilterCondition, value string) {
	return Equals, "HELLO WORLD"
}

func resolveCompoundFilterType(compoundFilterObject map[string]interface{}) CompoundFilterType {
	if compoundFilterObject["and"] != nil {
		return And
	} else {
		return Or
	}
}

func isPropertyField(propertyFilterObject map[string]interface{}) bool {
	return propertyFilterObject["property"] != nil
}
