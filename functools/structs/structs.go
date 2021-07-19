package structs

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func StructToMap(origin interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if origin == nil {
		return res
	}

	v := reflect.TypeOf(origin)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return res
	}

	reflectValue := reflect.ValueOf(origin)
	reflectValue = reflect.Indirect(reflectValue)

	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i)
		if tag != "" && tag != "-" {
			switch v.Field(i).Type.Kind() {
			case reflect.Struct:
				res[tag] = StructToMap(field.Interface())
			case reflect.Slice, reflect.Array:
				// 查询子元素类型
				kind := field.Index(0).Type().Kind()
				if kind == reflect.Ptr {
					kind = field.Index(0).Elem().Type().Kind()
				}
				switch kind {
				case reflect.Struct:
					val := []map[string]interface{}{}
					for i := 0; i < field.Len(); i++ {
						val = append(val, StructToMap(field.Index(i).Interface()))
					}
					res[tag] = val
				default:
					res[tag] = field.Interface()
				}
			default:
				res[tag] = field.Interface()
			}
		}
	}
	return res
}

func StructToValues(i interface{}, tag string) url.Values {
	v := url.Values{}
	if i != nil {
		iVal := reflect.ValueOf(i).Elem()
		tp := iVal.Type()
		for i := 0; i < iVal.NumField(); i++ {
			tag := tp.Field(i).Tag.Get(tag)
			if len(tag) > 0 {
				name := strings.Split(tag, ",")[0]
				if name != "-" {
					v.Set(name, fmt.Sprintf("%v", iVal.Field(i).Interface()))
				}
			}
		}
	}
	return v
}
