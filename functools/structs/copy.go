package structs

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

func CopyData(src interface{}, dest interface{}) error {
	data := StructToMap(src)
	return mapstructure.Decode(data, dest)
}

func CopyDataFromJson(src interface{}, dest interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}
