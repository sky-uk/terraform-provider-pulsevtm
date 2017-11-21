package util

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

// BuildStringArrayFromInterface : take an interface and convert it into an array of strings
func BuildStringArrayFromInterface(strings interface{}) []string {
	stringList := make([]string, len(strings.([]interface{})))
	for idx, stringValue := range strings.([]interface{}) {
		stringList[idx] = stringValue.(string)
	}
	return stringList
}

// BuildStringListFromSet : take an interface and convert it into an array of strings
func BuildStringListFromSet(strings *schema.Set) []string {
	stringList := make([]string, 0)
	for _, stringValue := range strings.List() {
		stringList = append(stringList, stringValue.(string))
	}
	return stringList
}

// AddSimpleGetAttributesToMap : wrapper for d.Get
func AddSimpleGetAttributesToMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributeNames []string) map[string]interface{} {

	for _, item := range attributeNames {
		attributeName := fmt.Sprintf("%s%s", attributeNamePrefix, item)
		attributeValue := d.Get(attributeName)
		switch attributeValue.(type) {
		case bool:
			mapItem[item] = attributeValue.(bool)
		case string:
			mapItem[item] = attributeValue.(string)
		case int:
			mapItem[item] = attributeValue.(int)
		case *schema.Set:
			mapItem[item] = attributeValue.(*schema.Set).List()
		default:
		}
	}
	return mapItem
}

// AddSimpleGetOkAttributesToMap : wrapper for d.GetOk
func AddSimpleGetOkAttributesToMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributeNames []string) map[string]interface{} {

	for _, item := range attributeNames {
		attributeName := fmt.Sprintf("%s%s", attributeNamePrefix, item)
		if attributeValue, ok := d.GetOk(attributeName); ok {
			switch attributeValue.(type) {
			case bool:
				mapItem[item] = attributeValue.(bool)
			case string:
				mapItem[item] = attributeValue.(string)
			case int:
				mapItem[item] = attributeValue.(int)
			default:
			}
		}
	}
	return mapItem
}

// AddChangedSimpleAttributesToMap : wrapper for d.HasChange & d.Get
func AddChangedSimpleAttributesToMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributeNames []string) map[string]interface{} {

	for _, item := range attributeNames {
		attributeName := fmt.Sprintf("%s%s", attributeNamePrefix, item)
		if d.HasChange(attributeName) {
			attributeValue := d.Get(attributeName)

			switch attributeValue.(type) {
			case bool:
				mapItem[item] = attributeValue.(bool)
			case string:
				mapItem[item] = attributeValue.(string)
			case int:
				mapItem[item] = attributeValue.(int)
			case *schema.Set:
				mapItem[item] = attributeValue.(*schema.Set).List()
			default:
			}
		}
	}
	return mapItem
}

// SetSimpleAttributesFromMap : wrapper for d.Set
func SetSimpleAttributesFromMap(d *schema.ResourceData, mapItem map[string]interface{}, attributeNamePrefix string, attributeNames []string) {

	for _, item := range attributeNames {
		attributeName := fmt.Sprintf("%s%s", attributeNamePrefix, item)
		d.Set(attributeName, mapItem[item])
	}
}

// BuildListMaps : builds a list of maps from a list of interfaces using a list of attributes
func BuildListMaps(itemList *schema.Set, attributeNames []string) ([]map[string]interface{}, error) {

	listOfMaps := make([]map[string]interface{}, 0)

	for _, item := range itemList.List() {

		definedItem := item.(map[string]interface{})
		newMap := make(map[string]interface{})

		for _, attributeName := range attributeNames {

			if attributeValue, ok := definedItem[attributeName]; ok {
				switch attributeValue.(type) {
				case bool:
					newMap[attributeName] = attributeValue.(bool)
				case string:
					newMap[attributeName] = attributeValue.(string)
				case int:
					newMap[attributeName] = attributeValue.(int)
				case []interface{}:
					newMap[attributeName] = attributeValue
				case *schema.Set:
					newMap[attributeName] = attributeValue.(*schema.Set).List()
				case map[string]interface{}:
					newMap[attributeName] = attributeValue.(map[string]interface{})
				case []map[string]interface{}:
					newMap[attributeName] = attributeValue.([]map[string]interface{})
				default:
					return listOfMaps, fmt.Errorf("util.BuildListMaps doesn't understand type for %+v", attributeValue)
				}
			}
		}
		listOfMaps = append(listOfMaps, newMap)
	}
	if len(listOfMaps) == 0 {
		emptyMap := make(map[string]interface{})
		listOfMaps = append(listOfMaps, emptyMap)
	}
	return listOfMaps, nil
}

// BuildReadListMaps : used by a read to build a list of maps which contain bools, strings, ints, float64s and lists of strings
func BuildReadListMaps(inputMap map[string]interface{}, attributeName string) (map[string]interface{}, error) {

	builtMap := make(map[string]interface{})

	for key, value := range inputMap {

		switch value.(type) {
		case bool:
			builtMap[key] = value.(bool)
		case string:
			builtMap[key] = value.(string)
		case float64:
			builtMap[key] = value.(float64)
		// []interface{} only configured / tested for a list of strings
		case []interface{}:
			builtMap[key] = schema.NewSet(schema.HashString, value.([]interface{}))
		default:
			return builtMap, fmt.Errorf("util.BuildReadListMaps doesn't understand type for %+v", value)
		}
	}
	return builtMap, nil
}
