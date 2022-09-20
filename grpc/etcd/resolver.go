package etcd

import (
	"context"

	"github.com/rs/xid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type etcdResolver struct {
	options *Options
	cli     *clientv3.Client
	em      endpoints.Manager
	key     string
}

func NewResolver() *etcdResolver {
	options := &Options{
		Id:  xid.New().String(),
		Url: "http://127.0.0.1:2379",
	}
	return &etcdResolver{
		options: options,
	}
}

func (r *etcdResolver) Init(opts ...Option) error {
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
	r.cli = cli
	r.em = em
	r.key = r.options.Service + "/" + r.options.Id
	return nil
}

func (r *etcdResolver) Register(addr string, leaseId clientv3.LeaseID, metadata interface{}) error {
	return r.em.AddEndpoint(context.Background(), r.key,
		endpoints.Endpoint{
			Addr:     addr,
			Metadata: metadata,
		},
		clientv3.WithLease(leaseId),
	)
}

func (r *etcdResolver) Unregister() error {
	return r.em.DeleteEndpoint(context.Background(), r.key)
}

func (r *etcdResolver) List() (endpoints.Key2EndpointMap, error) {
	return r.em.List(context.TODO())
}

func (r *etcdResolver) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	return r.cli.Grant(ctx, ttl)
}

func (r *etcdResolver) KeepAlive(ctx context.Context, id clientv3.LeaseID) error {
	respCh, err := r.cli.KeepAlive(ctx, id)
	if err != nil {
		return err
	}
	for {
		resp := <-respCh
		if resp == nil {
			return nil
		}
	}
}
