package protobuf

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type ProtobufJSON map[string]interface{}

func (p ProtobufJSON) GetTagMessage(tags ...string) []ProtobufJSON {
	defer func() {
		if err := recover(); err != nil {
			log.Warningf("GetTagMessage解析失败: %v", err)
		}
	}()
	output := []ProtobufJSON{}
	pJson := p
	for depth, tag := range tags {
		for _, v := range pJson.getOneLevelTagValue("message", tag) {
			value := v.(ProtobufJSON)
			if depth != len(tags)-1 {
				pJson = value
				break
			}
			output = append(output, value)
		}
	}
	return output
}

func (p ProtobufJSON) getOneLevelTagValue(wire string, tag string) []interface{} {
	vals := []interface{}{}
	for k, v := range p {
		if ok, _ := regexp.MatchString(fmt.Sprintf("%s:(\\d{2}):%s", tag, wire), k); ok {
			vals = append(vals, v)
		}
	}
	return vals
}

func (p ProtobufJSON) GetTagValue(wire string, tags ...string) []interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Warningf("GetTagValue解析失败: %v", err)
		}
	}()
	output := []interface{}{}
	pJson := p
	for depth, tag := range tags {
		levelWire := "message"
		if depth == len(tags)-1 {
			levelWire = wire
		}
		for _, v := range pJson.getOneLevelTagValue(levelWire, tag) {
			if depth != len(tags)-1 {
				pJson = v.(ProtobufJSON)
				break
			}
			var value interface{}
			switch wire {
			case "string":
				value = v.(string)
			case "varint", "fix32", "fix64":
				value = p.switchVarint(v)
			}
			output = append(output, value)
		}
	}
	return output
}

func (p ProtobufJSON) GetTagValueToString(tags ...string) string {
	values := p.GetTagValue("string", tags...)
	if len(values) > 0 {
		return values[0].(string)
	}
	return ""
}

func (p ProtobufJSON) GetTagValueToInt64(wire string, tags ...string) int64 {
	values := p.GetTagValue(wire, tags...)
	if len(values) > 0 {
		return values[0].(int64)
	}
	return 0
}

func (p ProtobufJSON) SetTagValue(tag string, value interface{}) {
	p[tag] = value
}

func (p ProtobufJSON) switchVarint(v interface{}) int64 {
	switch v.(type) {
	case int:
		return int64(v.(int))
	case int8:
		return int64(v.(int8))
	case int32:
		return int64(v.(int32))
	case int64:
		return int64(v.(int64))
	case uint:
		return int64(v.(uint))
	case uint8:
		return int64(v.(uint8))
	case uint16:
		return int64(v.(uint16))
	case uint32:
		return int64(v.(uint32))
	case uint64:
		return int64(v.(uint64))
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	}
	return 0
}
