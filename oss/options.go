package oss

import "context"

type Access struct {
	Endpoint string
	Key      string
	Secret   string
	Bucket   string
}

type Option func(o *Options)

type Options struct {
	Access
	Schema  string
	Context context.Context
}

func WithAccess(v Access) Option {
	return func(o *Options) {
		o.Access = v
	}
}

func WithSchema(schema string) Option {
	return func(o *Options) {
		o.Schema = schema
	}
}
