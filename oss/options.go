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
	TempBucket string
	Context    context.Context
}

func WithAccess(v Access) Option {
	return func(o *Options) {
		o.Access = v
	}
}

func WithTempBucket(temp string) Option {
	return func(o *Options) {
		o.TempBucket = temp
	}
}
