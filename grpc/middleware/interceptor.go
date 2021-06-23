package middleware

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	// set option
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	setLogLevel(options.logLevel)

	// return interceptor
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error

		startTime := time.Now().UnixNano() / 1e6

		// recover
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(ctx, r, options.recoveryFunc)
			}
		}()

		// 认证支持
		if options.authFunc != nil {
			ctx, err = authHandler(ctx, info, options.authFunc)
			if err != nil {
				return nil, err
			}
		}

		// 调用方法
		resp, err := handler(ctx, req)

		// 打印日志
		log.WithFields(log.Fields{
			"service":    info.FullMethod,
			"start_time": startTime,
			"time_used":  time.Now().UnixNano()/1e6 - startTime,
		}).Info("finished unary call")

		return resp, err
	}
}

// @TODO
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, stream)
	}
}
