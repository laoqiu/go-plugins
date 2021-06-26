package structs

import (
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type MyStruct struct {
	Name    string    `json:"name"`
	Address []Address `json:"address"`
	Photos  []string  `json:"photos"`
}

type Address struct {
	Addr string `json:"addr"`
}

func Test_CopyData(t *testing.T) {
	s := &MyStruct{Name: "test", Address: []Address{
		{
			Addr: "addr 1",
		},
		{
			Addr: "addr 2",
		},
	},
		Photos: []string{"photo1", "photo2"},
	}
	data := StructToMap(s)
	fmt.Println(data["address"].([]map[string]interface{})[0])
	fmt.Println(data)

	ss := &MyStruct{}
	mapstructure.Decode(data, ss)
	fmt.Println(ss, ss.Address[0])
}
