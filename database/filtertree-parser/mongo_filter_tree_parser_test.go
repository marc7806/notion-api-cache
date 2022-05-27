package filtertreeparser

import (
	"encoding/json"
	"testing"

	"github.com/marc7806/notion-cache/notion"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateNestedFilterTree(t *testing.T) {
	var inputData map[string]interface{}
	filterQueryJson := []byte(`
	{
        "and": [
            {
                "property": "Test Property 1",
                "text": {
                    "equals": "My custom value"
                }
            },
			{
                "property": "Test Property 2",
                "text": {
                    "equals": "My custom value 2"
                }
            }
        ]
    }
	`)
	err := json.Unmarshal(filterQueryJson, &inputData)
	if err != nil {
		t.Error("Error while unmarshalling input json string")
	}
	filterTree := notion.CreateFilterTree(inputData)
	mongoQuery := *ParseFilterTree(filterTree)

	if mongoQuery["$and"] == nil {
		t.Errorf("Wrong compound type. Expected '%s' but was '%s'", "$and", mongoQuery)
	}
	prop1 := *mongoQuery["$and"].([]interface{})[0].(*primitive.M)
	if prop1["properties.Test Property 1.value"].(primitive.M)["$eq"].(string) != "My custom value" {
		t.Errorf("Wrong property value. Expected '%s' but was '%s'", "My custom value", prop1["properties.Test Property 1.value"].(primitive.M)["$eq"])
	}
}
