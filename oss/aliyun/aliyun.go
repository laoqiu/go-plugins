package aliyun

import (
	"fmt"
	"io"
	"strings"

	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/laoqiu/go-plugins/oss"
)

type osssdk struct {
	client     *alioss.Client
	buckets    map[string]bool
	tempBucket string
	endpoint   string
}

func NewOSS() *osssdk {
	return &osssdk{
		buckets: make(map[string]bool),
	}
}

func (s *osssdk) Init(opts ...oss.Option) error {
	o := &oss.Options{}
	for _, opt := range opts {
		opt(o)
	}
	if o.TempBucket == "" {
		return fmt.Errorf("缺少temp bucket配置")
	}
	s.tempBucket = o.TempBucket
	// create client
	cli, err := alioss.New(o.Endpoint, o.Key, o.Secret)
	if err != nil {
		return err
	}
	s.client = cli
	s.endpoint = o.Access.Endpoint
	// get all buckets
	resp, err := s.client.ListBuckets()
	if err != nil {
		return err
	}
	for _, b := range resp.Buckets {
		s.buckets[b.Name] = true
	}
	if !s.buckets[s.tempBucket] {
		return fmt.Errorf("temp bucket not found")
	}
	return nil
}

func (s *osssdk) Upload(filename string, reader io.Reader) (string, error) {
	var out string
	if err := s.checkBucket(s.tempBucket); err != nil {
		return out, err
	}
	b, err := s.client.Bucket(s.tempBucket)
	if err != nil {
		return out, err
	}
	if err := b.PutObject(filename, reader); err != nil {
		return out, err
	}

	if strings.HasPrefix(filename, "/") {
		filename = filename[1:]
	}
	out = fmt.Sprintf("http://%s.%s%s", s.tempBucket, s.endpoint, filename)
	return out, nil
}

func (s *osssdk) Save(bucket, filename string) (string, error) {
	var out string
	if err := s.checkBucket(bucket); err != nil {
		return out, err
	}
	b, err := s.client.Bucket(bucket)
	if err != nil {
		return out, err
	}
	if _, err := b.CopyObjectFrom(s.tempBucket, filename, filename); err != nil {
		return out, err
	}
	if strings.HasPrefix(filename, "/") {
		filename = filename[1:]
	}
	out = fmt.Sprintf("http://%s.%s/%s", bucket, s.endpoint, filename)
	return out, nil
}

func (s *osssdk) checkBucket(bucket string) error {
	if !s.buckets[bucket] {
		return fmt.Errorf("bucket不存在")
	}
	return nil
}
