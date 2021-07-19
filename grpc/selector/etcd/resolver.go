package etcd

import (
	"context"

	"github.com/laoqiu/go-plugins/grpc/selector"
	"github.com/rs/xid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type etcdResolver struct {
	options *selector.Options
	em      endpoints.Manager
	key     string
}

func NewResolver() *etcdResolver {
	options := &selector.Options{
		Id:  xid.New().String(),
		Url: "http://127.0.0.1:2379",
	}
	return &etcdResolver{
		options: options,
	}
}

func (r *etcdResolver) Init(opts ...selector.Option) error {
	for _, o := range opts {
		o(r.options)
	}
	cli, err := clientv3.NewFromURL(r.options.Url)
	if err != nil {
		return err
	}
	em, err := endpoints.NewManager(cli, r.options.Service)
	if err != nil {
		return err
	}
	r.em = em
	r.key = r.options.Service + "/" + r.options.Id
	return nil
}

func (r *etcdResolver) Register(addr string, metadata interface{}) error {
	return r.em.AddEndpoint(context.Background(), r.key,
		endpoints.Endpoint{
			Addr:     addr,
			Metadata: metadata,
		},
	)
}

func (r *etcdResolver) Unregister() error {
	return r.em.DeleteEndpoint(context.Background(), r.key)
}
