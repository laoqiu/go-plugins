package selector

type Selector interface {
	Init(...Option) error
	Register(addr string, metadata interface{}) error
	Unregister() error
}
