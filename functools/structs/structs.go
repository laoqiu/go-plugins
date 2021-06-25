package structs

import (
	"fmt"
	"reflect"
	"strconv"
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
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

func MapToStruct(dest interface{}, data map[string]interface{}) error {
	v := reflect.TypeOf(dest)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("only accept struct or struct pointer; got %T", v)
	}

	reflectValue := reflect.ValueOf(dest)
	reflectValue = reflect.Indirect(reflectValue)

	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")

		if tag == "" || tag == "-" {
			continue
		}

		if val, ok := data[tag]; ok {

			var newVal reflect.Value
			fmt.Println(v.Field(i).Type.Kind())

			switch v.Field(i).Type.Kind() {
			case reflect.Struct, reflect.Ptr:
				var fieldDest interface{}

				field := reflectValue.Field(i)

				switch field.Kind() {
				case reflect.Ptr:
					if field.CanInterface() {
						fieldDest = field.Interface()
					}
				case reflect.Struct:
					if field.CanAddr() && field.CanInterface() {
						fieldDest = field.Addr().Interface()
					}
				}
				if fieldDest != nil {
					if err := MapToStruct(fieldDest, val.(map[string]interface{})); err != nil {
						return err
					}
				}
				return nil
			case reflect.ValueOf(val).Type().Kind():
				newVal = reflect.ValueOf(val)
			case reflect.String:
				newVal = reflect.ValueOf(fmt.Sprintf("%v", val))
			case reflect.Bool:
				newVal = reflect.ValueOf(asBool(val))
			case reflect.Float32:
				newVal = reflect.ValueOf(float32(asFloat(val)))
			case reflect.Float64:
				newVal = reflect.ValueOf(asFloat(val))
			case reflect.Int:
				newVal = reflect.ValueOf(int(asInt(val)))
			case reflect.Int8:
				newVal = reflect.ValueOf(int8(asInt(val)))
			case reflect.Int16:
				newVal = reflect.ValueOf(int16(asInt(val)))
			case reflect.Int32:
				newVal = reflect.ValueOf(int32(asInt(val)))
			case reflect.Int64:
				newVal = reflect.ValueOf(asInt(val))
			case reflect.Uint:
				newVal = reflect.ValueOf(uint(asInt(val)))
			case reflect.Uint8:
				newVal = reflect.ValueOf(uint8(asInt(val)))
			case reflect.Uint16:
				newVal = reflect.ValueOf(uint16(asInt(val)))
			case reflect.Uint32:
				newVal = reflect.ValueOf(uint32(asInt(val)))
			case reflect.Uint64:
				newVal = reflect.ValueOf(uint64(asInt(val)))
			default:
				return fmt.Errorf("no type match")
			}
			reflectValue.Field(i).Set(newVal)
		}
	}

	return nil
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func asFloat(src interface{}) float64 {
	switch v := src.(type) {
	case float64:
		return v
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.String:
		val, _ := strconv.ParseFloat(rv.String(), 64)
		return val
	case reflect.Float32:
		return rv.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Int())
	case reflect.Bool:
		if rv.Bool() {
			return 1
		}
		return 0
	}
	return 0
}

func asInt(src interface{}) int64 {
	switch v := src.(type) {
	case int64:
		return v
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.String:
		val, _ := strconv.ParseInt(rv.String(), 10, 61)
		return val
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(rv.Uint())
	case reflect.Bool:
		if rv.Bool() {
			return 1
		}
		return 0
	}
	return 0
}

func asBool(src interface{}) bool {
	switch v := src.(type) {
	case bool:
		return v
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.String:
		if rv.String() == "true" {
			return true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv.Int() > 0 {
			return true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rv.Uint() > 0 {
			return true
		}
	}
	return false
}
