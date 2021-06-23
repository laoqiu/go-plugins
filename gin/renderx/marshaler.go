package renderx

import (
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

type ResponseWithRequestID struct {
	id     int
	err    error
	result proto.Message
}

func NewResponseWithRequestID(requestID int, result proto.Message, err error) *ResponseWithRequestID {
	return &ResponseWithRequestID{id: requestID, result: result, err: err}
}

func batchUnmarshal(resp []interface{}) (interface{}, error) {
	result := []interface{}{}

	for _, data := range resp {
		if obj, ok := data.(*ResponseWithRequestID); ok {
			dataMap := map[string]interface{}{
				"request_id": obj.id,
			}
			if obj.err != nil {
				dataMap["error"] = obj.err.Error()
			}
			if obj.result != nil {
				dataMap["result"], _ = protobufToMap(obj.result)
			}
			result = append(result, dataMap)
		}
	}

	return result, nil
}

func protobufToMap(data interface{}) (map[string]interface{}, error) {
	dataMap := make(map[string]interface{})

	// 处理protobuf数据
	if data != nil {
		m := jsonpb.Marshaler{EmitDefaults: true, OrigName: true}
		data, err := m.MarshalToString(data.(proto.Message))
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(data), &dataMap); err != nil {
			return nil, err
		}
	}
	return dataMap, nil
}
