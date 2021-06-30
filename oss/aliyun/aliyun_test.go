package aliyun

import (
	"fmt"
	"os"
	"testing"

	"github.com/laoqiu/go-plugins/oss"
)

func Test_Aliyun(t *testing.T) {
	sdk := NewOSS()
	if err := sdk.Init(
		oss.WithAccess(oss.Access{
			Endpoint: "oss-cn-hangzhou.aliyuncs.com",
			Key:      "",
			Secret:   "",
			Bucket:   "bcltemp",
		})); err != nil {
		t.Error(err)
		return
	}

	file, err := os.Open("1.gif")
	if err != nil {
		t.Error(err)
		return
	}

	output, err := sdk.Upload("temp/test.gif", file)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(output)

	newPath, err := sdk.Save(output, "box/test.gif")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(newPath)
}
