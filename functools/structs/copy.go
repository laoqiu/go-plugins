package structs

import (
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

func CopyData(src interface{}, dest interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}

func CopyProtobufData(src proto.Message, desc interface{}) error {
	m := jsonpb.Marshaler{EmitDefaults: true, OrigName: true}
	data, err := m.MarshalToString(src)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(data), desc); err != nil {
		return err
	}
	return nil
}
