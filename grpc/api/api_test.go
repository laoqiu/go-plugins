package api

import (
	"context"
	"testing"
	"time"
)

func Test_Run(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	if err := Run(ctx, &Config{
		Port:       ":8081",
		EtcdServer: "127.0.0.1:2379",
		Services: []string{
			"micro.financial",
		},
	}); err != nil {
		t.Error(err)
	}
}
