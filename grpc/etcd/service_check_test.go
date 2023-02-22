package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

func TestServiceCheck(t *testing.T) {

	// 待检测的服务
	serviceList := []string{"micro.oss"}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cli, err := clientv3.NewFromURL("10.9.6.232:2379")
	if err != nil {
		t.Error(err)
		return
	}

	managerMap := map[string]endpoints.Manager{}

	for _, serviceName := range serviceList {
		em, err := endpoints.NewManager(cli, serviceName)
		if err != nil {
			t.Error(err)
			return
		}
		managerMap[serviceName] = em
	}

	var lastErrorPushTime int64
	timer := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			fmt.Println("in timer", time.Now().Unix())
			for k, v := range managerMap {
				endpoints, err := v.List(context.TODO())
				if err != nil {
					continue
				}
				if len(endpoints) == 0 {
					fmt.Printf("服务 %s 已失联\n", k)
					if (time.Now().Unix() - lastErrorPushTime) > 3600 {
						lastErrorPushTime = time.Now().Unix()
						// TODO 发送webhook
					}
				}
			}
		}
	}

}
