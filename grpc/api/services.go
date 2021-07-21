package api

import (
	"context"
	"fmt"
	"strings"
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
	cc      grpc.ClientConnInterface
}

func ListServices(addr string, serviceNames []string) (map[string]*Service, error) {
	cli, err := clientv3.NewFromURL(addr)
	if err != nil {
		return nil, err
	}

	serviceMap := make(map[string]*Service)
	for _, sName := range serviceNames {
		serviceMap[sName] = &Service{
			name: sName,
		}
		etcdResolver, err := resolver.NewBuilder(cli)
		if err != nil {
			return nil, err
		}
		conn, err := grpc.Dial("etcd:///"+sName,
			grpc.WithResolvers(etcdResolver),
			// grpc.WithBalancerName("round_robin"),
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
			grpc.WithInsecure(),
		)
		if err != nil {
			return nil, err
		}
		serviceMap[sName].cc = conn

		// 勿必增加超时处理
		// 如果不存在服务会返回访问超时
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		methods, err := getMethods(getDescSource(ctx, conn))
		if err != nil {
			return nil, err
		}
		serviceMap[sName].methods = methods
	}

	return serviceMap, nil
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
		for _, md := range sd.GetMethods() {
			descs = append(descs, md)
		}
	}

	return descs, nil
}

func splitMethodName(name string) (svc, method string) {
	dot := strings.LastIndex(name, ".")
	slash := strings.LastIndex(name, "/")
	sep := dot
	if slash > dot {
		sep = slash
	}
	if sep < 0 {
		return "", name
	}
	return name[:sep], name[sep+1:]
}
