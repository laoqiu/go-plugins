package protobuf

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"
)

// MS4wLjABAAAAyGEg-noNnFOka6GHED6TdiOktn1iX6pLFKiXtBOkyh8CjKkY2rvBFR4GleqVspHm
// MS4wLjABAAAAyGEg-noNnFOka6GHED6TdiOktn1iX6pLFKiXtBOkyh8CjKkY2rvBFR4GleqVspHm

func Test_Buffer(t *testing.T) {
	// s := ""
	sb, _ := ioutil.ReadFile("init.txt")
	// b, err := ioutil.ReadFile("/Users/laoqiu/Downloads/get_by_user_init")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	b, err := base64.StdEncoding.DecodeString(string(sb))
	if err != nil {
		t.Error(err)
	}
	buf := NewBuffer(b)
	if pjson, ok := buf.Unmarshal(); ok {
		// fmt.Println(pjson)
		// o, _ := json.Marshal(pjson)
		// ioutil.WriteFile("output.txt", o, 0666)
		// pjson.GetTagMessage("08")
		messages := pjson.GetTagMessage("06", "203", "02")
		fmt.Println(len(messages))
		for _, m := range messages {
			relationId := m.GetTagValueToString("01")
			messageId := m.GetTagValueToInt64("varint", "02")
			sessionId := m.GetTagValueToString("04")
			fmt.Println(relationId, messageId, sessionId)
		}
		// j, _ := json.Marshal(pjson)
		// fmt.Println(string(j))
		// fmt.Println(pjson.GetTagValueToInt64("varint", "13"))

	}
}
