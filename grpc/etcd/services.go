package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type Service struct {
	name    string
	methods []*desc.MethodDescriptor
	cc      *grpc.ClientConn
}

func FindService(cli *clientv3.Client, name string) (*Service, error) {
	service := &Service{
		name: name,
	}

	etcdResolver, err := resolver.NewBuilder(cli)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial("etcd:///"+name,
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	service.cc = conn

	// 勿必增加超时处理
	// 如果不存在服务会返回访问超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	methods, err := getMethods(getDescSource(ctx, conn))
	if err != nil {
		return nil, err
	}
	service.methods = methods

	return service, nil
}

func getDescSource(ctx context.Context, cc grpc.ClientConnInterface, reflHeaders ...string) grpcurl.DescriptorSource {
	md := grpcurl.MetadataFromHeaders(reflHeaders)
	refCtx := metadata.NewOutgoingContext(ctx, md)
	refClient := grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(cc))
	return grpcurl.DescriptorSourceFromServer(ctx, refClient)
}

func getMethods(source grpcurl.DescriptorSource) ([]*desc.MethodDescriptor, error) {
	allServices, err := source.ListServices()
	if err != nil {
		return nil, err
	}

	var descs []*desc.MethodDescriptor
	for _, svc := range allServices {
		if svc == "grpc.reflection.v1alpha.ServerReflection" {
			continue
		}
		d, err := source.FindSymbol(svc)
		if err != nil {
			return nil, err
		}
		sd, ok := d.(*desc.ServiceDescriptor)
		if !ok {
			return nil, fmt.Errorf("%s should be a service descriptor but instead is a %T", d.GetFullyQualifiedName(), d)
		}
		descs = append(descs, sd.GetMethods()...)
	}

	return descs, nil
}
