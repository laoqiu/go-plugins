package api

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port       string   `json:"port"`
	EtcdServer string   `json:"etcd_server"`
	Services   []string `json:"services"`
}

func Run(ctx context.Context, conf *Config) error {
	services, err := ListServices(conf.EtcdServer, conf.Services)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/invoke", RPCInvokeHandler(services))

	srv := &http.Server{
		Addr:    conf.Port,
		Handler: mux,
	}

	log.Infof("grpc api server at %s", conf.Port)
	defer log.Info("grpc api server exit")

	go func() {
		<-ctx.Done()
		srv.Close()
	}()

	return srv.ListenAndServe()
}
