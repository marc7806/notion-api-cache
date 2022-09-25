package sortparser

import (
	"encoding/json"
	"testing"

	"github.com/marc7806/notion-cache/types"
)

type TestBody struct {
	Sorts []types.QuerySort `json:"sorts"`
}

func TestCreateNestedFilterTree(t *testing.T) {
	var inputData TestBody
	sortQueryJson := []byte(`
	{
		"sorts": [
			{
				"property": "Test Property 1",
				"direction": "ascending"
			},
			{
				"property": "Test Property 2",
				"direction": "descending"
			}
		]
    }
	`)
	err := json.Unmarshal(sortQueryJson, &inputData)
	if err != nil {
		t.Error("Error while unmarshalling input json string")
	}
	mongoQuery := *ParseSortOptions(inputData.Sorts)

	if mongoQuery["properties.Test Property 1.value"] == nil {
		t.Errorf("Wrong sort type. Expected '%s' but was '%s'", "Test Property 1", mongoQuery)
	}
	prop1 := mongoQuery["properties.Test Property 1.value"].(int8)
	if prop1 != 1 {
		t.Errorf("Wrong sort value. Expected '%d' but was '%d'", 1, prop1)
	}

	if mongoQuery["properties.Test Property 2.value"] == nil {
		t.Errorf("Wrong sort type. Expected '%s' but was '%s'", "Test Property 2", mongoQuery)
	}
	prop2 := mongoQuery["properties.Test Property 2.value"].(int8)
	if prop2 != -1 {
		t.Errorf("Wrong sort value. Expected '%d' but was '%d'", -1, prop1)
	}
}
