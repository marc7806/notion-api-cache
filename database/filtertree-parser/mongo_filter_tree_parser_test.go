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

func TestDoesNotContainFilter(t *testing.T) {
	var inputData map[string]interface{}
	filterQueryJson := []byte(`
	{
        "property": "Test Property 1",
        "text": {
            "does_not_contain": "vorläuf. Freigabe Live-CC"
        }
    }
	`)
	err := json.Unmarshal(filterQueryJson, &inputData)
	if err != nil {
		t.Error("Error while unmarshalling input json string")
	}
	filterTree := notion.CreateFilterTree(inputData)
	mongoQuery := *ParseFilterTree(filterTree)

	condition, ok := mongoQuery["properties.Test Property 1.value"].(primitive.M)
	if !ok {
		t.Fatalf("Missing condition for property. Query was '%s'", mongoQuery)
	}
	regex, ok := condition["$not"].(primitive.Regex)
	if !ok {
		t.Fatalf("Wrong condition. Expected '$not' with regex value but was '%s'", condition)
	}
	// the dot is a regex metacharacter and has to be escaped to match literally
	expectedPattern := `vorläuf\. Freigabe Live-CC`
	if regex.Pattern != expectedPattern {
		t.Errorf("Wrong regex pattern. Expected '%s' but was '%s'", expectedPattern, regex.Pattern)
	}
}
