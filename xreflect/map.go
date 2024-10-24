package xreflect

import (
	"reflect"
	"strings"
)

func ToMap(obj any, tagName string) map[any]any {
	result := make(map[any]any)

	objType := reflect.TypeOf(obj)
	objVal := reflect.ValueOf(obj)

	if objType.Kind() == reflect.Pointer {
		objType = objType.Elem()
		objVal = objVal.Elem()
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldName := field.Name

		tagValue := field.Tag.Get(tagName)
		tagValue, _, _ = strings.Cut(tagValue, ",")
		if tagValue == "" {
			continue
		}
		result[tagValue] = objVal.FieldByName(fieldName).Interface()
	}

	return result
}

func ToMapString(obj any, tagName string) map[string]any {
	result := make(map[string]any)

	objType := reflect.TypeOf(obj)
	objVal := reflect.ValueOf(obj)

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldName := field.Name

		tagValue := field.Tag.Get(tagName)
		tagValue, _, _ = strings.Cut(tagValue, ",")
		if tagValue == "" {
			continue
		}
		result[tagValue] = objVal.FieldByName(fieldName).Interface()
	}

	return result
}
