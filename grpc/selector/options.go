package selector

type Option func(o *Options)

type Options struct {
	Id      string
	Url     string
	Service string
}

func WithURL(url string) Option {
	return func(o *Options) {
		o.Url = url
	}
}

func WithService(service string) Option {
	return func(o *Options) {
		o.Service = service
	}
}
