package notion

import (
	"errors"
	"log"
)

type fn func(map[string]interface{}, string) string

var (
	resolvers map[string]fn
)

func init() {
	log.Println("Initializing Notion property type resolvers")
	resolvers = map[string]fn{
		"rich_text":        resolvePlainText,
		"title":            resolvePlainText,
		"last_edited_time": resolveLastEditedTime,
		"created_time":     resolveCreatedTime,
		"select":           resolveSelect,
		"multi_select":     resolveMultiSelect,
		"formula":          resolveFormula,
	}
}

func ResolvePropertyType(property map[string]interface{}) (propValue string, err error) {
	propType := property["type"].(string)
	resolver, ok := resolvers[propType]
	if ok {
		return resolver(property, propType), nil
	}
	return "", errors.New("no matching resolver found")
}

func resolvePlainText(prop map[string]interface{}, propType string) string {
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

func resolveLastEditedTime(prop map[string]interface{}, propType string) string {
	return prop[propType].(string)
}

func resolveCreatedTime(prop map[string]interface{}, propType string) string {
	return prop[propType].(string)
}

func resolveSelect(prop map[string]interface{}, propType string) string {
	selectData := prop[propType]
	if selectData == nil {
		return ""
	}
	return selectData.(map[string]interface{})["name"].(string)
}

func resolveMultiSelect(prop map[string]interface{}, propType string) string {
	multiSelectOptions := prop[propType].([]interface{})
	result := ""

	if len(multiSelectOptions) == 0 {
		return result
	}

	for _, dataEntry := range multiSelectOptions {
		// currently we concat multi values with ; seperator to a string
		result += dataEntry.(map[string]interface{})["name"].(string) + ";"
	}
	return result
}

func resolveFormula(prop map[string]interface{}, propType string) string {
	return prop[propType].(map[string]interface{})["string"].(string)
}
