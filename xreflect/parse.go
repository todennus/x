package xreflect

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/todennus/x/conversion"
)

func Parse(obj any, strict bool, tagName string, fieldVal func(string) any) error {
	objType := reflect.TypeOf(obj).Elem()
	objVal := reflect.ValueOf(obj).Elem()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldName := field.Name

		tagValue := field.Tag.Get(tagName)
		tagValue, _, _ = strings.Cut(tagValue, ",")
		if tagValue == "" {
			continue
		}

		fieldValue := fieldVal(tagValue)
		if fieldValue == "" {
			continue
		}

		fieldVal := objVal.FieldByName(fieldName)

		switch field.Type.Kind() {
		case reflect.String:
			s, err := conversion.ToString(fieldValue, strict)
			if err != nil {
				return fmt.Errorf("%s: %s", tagValue, err.Error())
			}

			fieldVal.SetString(s)

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			intFormVal, err := conversion.ToInt(fieldValue, strict)
			if err != nil {
				return fmt.Errorf("%s: %s", tagValue, err.Error())
			}
			fieldVal.SetInt(intFormVal)

		case reflect.Float32, reflect.Float64:
			floatFormVal, err := conversion.ToFloat(fieldValue, strict)
			if err != nil {
				return fmt.Errorf("%s: %s", tagValue, err.Error())
			}
			fieldVal.SetFloat(floatFormVal)

		case reflect.Bool:
			boolFormVal, err := conversion.ToBool(fieldValue, strict)
			if err != nil {
				return fmt.Errorf("%s: %s", tagValue, err.Error())
			}
			fieldVal.SetBool(boolFormVal)

		default:
			return fmt.Errorf("not support type %s", field.Type.Kind())
		}
	}

	return nil
}
