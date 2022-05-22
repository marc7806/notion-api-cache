package notion

import (
	"fmt"
)

type RichText struct {
	Type        string      `json:"type"`
	Text        interface{} `json:"text"`
	Annotations interface{} `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        interface{} `json:"href"`
}

type fn func(map[string]interface{}, string) interface{}

var (
	resolvers map[string]fn
)

func init() {
	fmt.Println("Initializinggg")
	resolvers = map[string]fn{
		"rich_text": resolveRichText,
	}
}

func ResolvePropertyType(property map[string]interface{}) interface{} {
	propType := property["type"].(string)
	resolver, ok := resolvers[propType]
	if ok {
		return resolver(property, propType)
	}
	return "EMPTY"
}

func resolveRichText(prop map[string]interface{}, propType string) interface{} {
	richTextData := prop[propType].([]interface{})
	result := ""

	if len(richTextData) == 0 {
		return result
	}

	for _, dataEntry := range richTextData {
		result += dataEntry.(map[string]interface{})["plain_text"].(string)
	}
	return result
}
