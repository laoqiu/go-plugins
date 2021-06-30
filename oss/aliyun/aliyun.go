package aliyun

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/laoqiu/go-plugins/oss"
)

var DefaultSchema = "http"

type osssdk struct {
	client   *alioss.Client
	bucket   string
	endpoint string
	schema   string
}

func NewOSS() *osssdk {
	return &osssdk{}
}

func (s *osssdk) Init(opts ...oss.Option) error {
	o := &oss.Options{
		Schema: DefaultSchema,
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.Bucket == "" {
		return fmt.Errorf("缺少bucket配置")
	}
	// create client
	cli, err := alioss.New(o.Endpoint, o.Key, o.Secret)
	if err != nil {
		return err
	}
	s.client = cli
	s.endpoint = o.Access.Endpoint
	s.bucket = o.Bucket
	s.schema = o.Schema
	return nil
}

func (s *osssdk) Upload(filename string, reader io.Reader) (string, error) {
	var out string
	b, err := s.client.Bucket(s.bucket)
	if err != nil {
		return out, err
	}
	if err := b.PutObject(filename, reader); err != nil {
		return out, err
	}

	return s.toAbsolutePath(filename), nil
}

func (s *osssdk) Save(from, to string) (string, error) {
	// 如果from是绝对路径，转成相对路径
	if strings.HasPrefix(from, "http") {
		u, err := url.Parse(from)
		if err != nil {
			return "", fmt.Errorf("不是有效的路径参数[from]")
		}
		from = u.Path[1:]
	}
	var out string
	b, err := s.client.Bucket(s.bucket)
	if err != nil {
		return out, err
	}
	if _, err := b.CopyObject(from, to); err != nil {
		return out, fmt.Errorf("复制文件出错: %v", err)
	}
	return s.toAbsolutePath(to), nil
}

func (s *osssdk) toAbsolutePath(path string) string {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	return fmt.Sprintf("%s://%s.%s/%s", s.schema, s.bucket, s.endpoint, path)
}
