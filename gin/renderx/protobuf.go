package renderx

import (
	"encoding/json"
	"net/http"
)

// ProtoBuf contains the given interface object.
type ProtoBuf struct {
	Data interface{}
}

// var protobufContentType = []string{"application/x-protobuf"}
var jsonContentType = []string{"application/json; charset=utf-8"}

// Render (ProtoBuf) marshals the given interface object and writes data with custom ContentType.
func (r ProtoBuf) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	result, err := protobufToMap(r.Data)
	if err != nil {
		return err
	}
	// 自定义外层数据格式
	jsonBytes, err := json.Marshal(Success(result))
	if err != nil {
		return err
	}

	_, err = w.Write(jsonBytes)
	return err
}

func (r ProtoBuf) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = jsonContentType
	}
}

// 批处理protobuf数据
type BatchProtobuf struct {
	Data []interface{}
}

func (r BatchProtobuf) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	result, err := batchUnmarshal(r.Data)
	if err != nil {
		return err
	}
	// 自定义外层数据格式
	jsonBytes, err := json.Marshal(Success(result))
	if err != nil {
		return err
	}

	_, err = w.Write(jsonBytes)
	return err
}

func (r BatchProtobuf) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = jsonContentType
	}
}
