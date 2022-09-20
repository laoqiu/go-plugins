package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/laoqiu/go-plugins/grpc/etcd"
	"github.com/laoqiu/go-plugins/test/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type testHandler struct {
	proto.UnimplementedTestServerServer
}

func (h *testHandler) Ping(ctx context.Context, req *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{Id: "test"}, nil
}

func main() {
	srv := grpc.NewServer()

	c, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	proto.RegisterTestServerServer(srv, &testHandler{})

	reflection.Register(srv)

	r := etcd.NewResolver()
	if err := r.Init(etcd.WithService("micro.test"), etcd.WithURL("10.9.6.232:2379")); err != nil {
		log.Fatalln(err)
		return
	}
	lease, _ := r.Grant(context.TODO(), 3)
	if err := r.Register("127.0.0.1:50051", lease.ID, nil); err != nil {
		log.Fatalln(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		defer fmt.Println("keep alive exit")
		respCh, err := r.KeepAlive(ctx, lease.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			resp := <-respCh
			if resp == nil {
				break
			}
		}
	}(ctx)

	fmt.Println("start server")

	srv.Serve(c)

	fmt.Println("stop server")
}
