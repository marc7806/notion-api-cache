package notion

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

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
	And CompoundFilterType = "and"
	Or                     = "or"
)

type FilterCondition string

const (
	Equals         FilterCondition = "equals"
	DoesNotEqual                   = "does_not_equal"
	Contains                       = "contains"
	DoesNotContain                 = "does_not_contain"
	StartsWith                     = "starts_with"
	EndsWith                       = "ends_with"
	IsEmpty                        = "is_empty"
	IsNotEmpty                     = "is_not_empty"
)

func CreateFilterTree(filterQuery map[string]interface{}) *FilterTreeNode {
	// recursive parse query body from rest endpoint using dfs and return filter tree node
	return traverseFilterQuery(filterQuery)
}

func traverseFilterQuery(filterQuery map[string]interface{}) *FilterTreeNode {
	// 1. case: single property filter - Basecase
	if isPropertyField(filterQuery) {
		condition, value, err := resolveFilterCondition(filterQuery)
		if err != nil {
			return nil
		}
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
	filterTreeNode := &FilterTreeNode{
		CompoundType: &filterCompoundType,
	}
	for _, el := range filterQuery[string(filterCompoundType)].([]interface{}) {
		filterTreeNode.Children = append(filterTreeNode.Children, *traverseFilterQuery(el.(map[string]interface{})))
	}
	return filterTreeNode
}

func resolveFilterCondition(propertyFilterObject map[string]interface{}) (condition FilterCondition, value string, err error) {
	for key, val := range propertyFilterObject {
		if key != "property" {
			// here we have our filter condition
			for conditionKey, conditionVal := range val.(map[string]interface{}) {
				convertedVal := ""
				if reflect.TypeOf(conditionVal).Kind() == reflect.Bool {
					convertedVal = strconv.FormatBool(conditionVal.(bool))
				} else {
					convertedVal = conditionVal.(string)
				}
				return FilterCondition(conditionKey), convertedVal, nil
			}
		}
	}
	propString, _ := json.Marshal(propertyFilterObject)
	return Equals, "", errors.New("Error: Query contains not supported filter condition: " + string(propString))
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
