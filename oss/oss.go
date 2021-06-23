package oss

import (
	"io"
)

type OSS interface {
	// 初始化
	Init(...Option) error
	// 上传到临时目录
	Upload(filename string, reader io.Reader) (string, error)
	// 保存到正式目录
	Save(bucket, fullname string) (string, error)
}
