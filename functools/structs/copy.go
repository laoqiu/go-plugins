package structs

import "encoding/json"

func CopyData(src interface{}, dest interface{}) error {
	return MapToStruct(dest, StructToMap(src))
}

func CopyDataFromJson(src interface{}, dest interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}
