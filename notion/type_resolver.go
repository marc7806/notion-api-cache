package notion

import "log"

type fn func(map[string]interface{}, string) interface{}

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

func ResolvePropertyType(property map[string]interface{}) interface{} {
	propType := property["type"].(string)
	resolver, ok := resolvers[propType]
	if ok {
		return resolver(property, propType)
	}
	return nil
}

func resolvePlainText(prop map[string]interface{}, propType string) interface{} {
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

func resolveLastEditedTime(prop map[string]interface{}, propType string) interface{} {
	return prop[propType]
}

func resolveCreatedTime(prop map[string]interface{}, propType string) interface{} {
	return prop[propType]
}

func resolveSelect(prop map[string]interface{}, propType string) interface{} {
	selectData := prop[propType]
	if selectData == nil {
		return ""
	}
	return selectData.(map[string]interface{})["name"]
}

func resolveMultiSelect(prop map[string]interface{}, propType string) interface{} {
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

func resolveFormula(prop map[string]interface{}, propType string) interface{} {
	return prop[propType].(map[string]interface{})["string"]
}
