package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	test "github.com/laoqiu/go-plugins/test/proto"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
)

func TestResolver(t *testing.T) {
	r := NewResolver()
	if err := r.Init(WithService("micro.test"), WithURL("10.9.6.232:2379")); err != nil {
		t.Error(err)
		return
	}

	etcdResolver, err := resolver.NewBuilder(r.cli)
	if err != nil {
		t.Error(err)
	}
	conn, _ := grpc.Dial("etcd:///"+"micro.test", grpc.WithResolvers(etcdResolver), grpc.WithInsecure())
	c := test.NewTestServerClient(conn)

	count := 0
	timer := time.NewTicker(time.Second * 2)
	for {
		<-timer.C
		fmt.Println(r.List())
		resp, err := c.Ping(context.Background(), &test.PingRequest{Id: "1111"})
		if err != nil {
			fmt.Println("err: ", err)
		} else {
			fmt.Println("response", resp.Id)
		}
		count += 1
		if count == 12 {
			break
		}
	}
}
