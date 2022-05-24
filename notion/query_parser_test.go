package notion

import (
	"encoding/json"
	"testing"
)

func TestCreateNestedFilterTree(t *testing.T) {
	var inputData map[string]interface{}
	filterQueryJson := []byte(`
	{
        "and": [
            {
                "property": "My Property 1",
                "text": {
                    "is_not_empty": true
                }
            },
            {
                "property": "My Property 2",
                "text": {
                    "equals": "Test Value"
                }
            }
        ]
    }
	`)
	err := json.Unmarshal(filterQueryJson, &inputData)
	if err != nil {
		t.Error("Error while unmarshalling input json string")
	}

	filterTree := CreateFilterTree(inputData)

	if *filterTree.CompoundType != And {
		t.Errorf("Wrong compound type. Expected '%s' but was '%s'", And, *filterTree.CompoundType)
	}
	if len(filterTree.Children) != 2 {
		t.Errorf("Wrong children size. Expected %d but was %d", 2, len(filterTree.Children))
	}
	if filterTree.Children[0].Operation.Property != "My Property 1" {
		t.Errorf("Wrong child property name. Expected '%s' but was '%s'", "My Property 1", filterTree.Children[0].Operation.Property)
	}
	if filterTree.Children[0].Operation.Condition != IsNotEmpty {
		t.Errorf("Wrong child operation condition. Expected '%s' but was '%s'", IsNotEmpty, filterTree.Children[0].Operation.Condition)
	}
	if filterTree.Children[1].Operation.Value != "Test Value" {
		t.Errorf("Wrong child operation value. Expected '%s' but was '%s'", "Test Value", filterTree.Children[1].Operation.Value)
	}
}

func TestCreateSingleFilterTree(t *testing.T) {
	var inputData map[string]interface{}
	filterQueryJson := []byte(`
	{
		"property": "My Property 1",
		"text": {
			"is_not_empty": true
		}
	}
	`)
	err := json.Unmarshal(filterQueryJson, &inputData)
	if err != nil {
		t.Error("Error while unmarshalling input json string")
	}

	filterTree := CreateFilterTree(inputData)

	if filterTree.Operation.Property != "My Property 1" {
		t.Errorf("Wrong child property name. Expected '%s' but was '%s'", "My Property 1", filterTree.Children[0].Operation.Property)
	}
	if filterTree.Operation.Condition != IsNotEmpty {
		t.Errorf("Wrong child operation condition. Expected '%s' but was '%s'", IsNotEmpty, filterTree.Children[0].Operation.Condition)
	}
	if filterTree.Operation.Value != "true" {
		t.Errorf("Wrong child operation value. Expected '%s' but was '%s'", "true", filterTree.Children[1].Operation.Value)
	}
}

func TestCreateDeepNestedFilterTree(t *testing.T) {
	var inputData map[string]interface{}
	filterQueryJson := []byte(`
	{
        "and": [
            {
                "property": "My Property 1",
                "text": {
                    "is_not_empty": true
                }
            },
            {
				"or": [
					{
						"property": "My Property 3",
						"text": {
							"equals": "Test Value"
						}
					},
					{
						"property": "My Property 4",
						"text": {
							"equals": "Test Value"
						}
					}
				]
			}
        ]
    }
	`)
	err := json.Unmarshal(filterQueryJson, &inputData)
	if err != nil {
		t.Error("Error while unmarshalling input json string")
	}

	filterTree := CreateFilterTree(inputData)

	if *filterTree.CompoundType != And {
		t.Errorf("Wrong compound type. Expected '%s' but was '%s'", And, *filterTree.CompoundType)
	}
	if *filterTree.Children[1].CompoundType != Or {
		t.Errorf("Wrong compound type. Expected '%s' but was '%s'", Or, *filterTree.Children[0].CompoundType)
	}
	if len(filterTree.Children[1].Children) != 2 {
		t.Errorf("Wrong children size. Expected %d but was %d", 2, len(filterTree.Children[0].Children))
	}
	if filterTree.Children[0].Operation.Property != "My Property 1" {
		t.Errorf("Wrong child property name. Expected '%s' but was '%s'", "My Property 1", filterTree.Children[0].Operation.Property)
	}
	if filterTree.Children[0].Operation.Condition != IsNotEmpty {
		t.Errorf("Wrong child operation condition. Expected '%s' but was '%s'", IsNotEmpty, filterTree.Children[0].Operation.Condition)
	}

	if filterTree.Children[1].Children[0].Operation.Value != "Test Value" {
		t.Errorf("Wrong child operation value. Expected '%s' but was '%s'", "Test Value", filterTree.Children[1].Children[0].Operation.Value)
	}
	if filterTree.Children[1].Children[0].Operation.Condition != Equals {
		t.Errorf("Wrong child operation value. Expected '%s' but was '%s'", Equals, filterTree.Children[1].Children[0].Operation.Condition)
	}
}
