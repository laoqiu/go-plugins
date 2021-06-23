package middleware

import (
	"context"
	"fmt"
	"net"
	"testing"

	pb "github.com/laoqiu/go-plugins/test/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type testHandler struct {
	pb.UnimplementedTestServerServer
}

func NewTestHandler() pb.TestServerServer {
	return &testHandler{}
}

func (h *testHandler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	// panic("nnnnnnn")
	return &pb.PingResponse{}, nil
}

func init() {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(
			WithLogLevel(logrus.DebugLevel),
			WithRecoverFunc(func(ctx context.Context, v interface{}) error {
				fmt.Println("WithRecoverFunc =>", v)
				return nil
			}),
			WithAuthFunc(func(ctx context.Context) (context.Context, error) {
				fmt.Println("WithAuthFunc =>", ctx)
				return ctx, nil
			}),
		)),
	)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	pb.RegisterTestServerServer(server, NewTestHandler())
	go server.Serve(lis)
}

func Test_UnaryServerInterceptor(t *testing.T) {
	conn, _ := grpc.Dial(":50051", grpc.WithInsecure())
	cli := pb.NewTestServerClient(conn)
	resp, err := cli.Ping(context.Background(), &pb.PingRequest{Id: "1111"})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(resp.Id)
}
