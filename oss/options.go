package oss

import "context"

type Access struct {
	Endpoint string
	Key      string
	Secret   string
}

type Option func(o *Options)

type Options struct {
	Access
	Bucket  string
	Schema  string
	Context context.Context
}

func WithAccess(v Access) Option {
	return func(o *Options) {
		o.Access = v
	}
}

func WithBucket(temp string) Option {
	return func(o *Options) {
		o.Bucket = temp
	}
}

func WithSchema(schema string) Option {
	return func(o *Options) {
		o.Schema = schema
	}
}
