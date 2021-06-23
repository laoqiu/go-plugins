package middleware

import (
	"context"

	"google.golang.org/grpc"
)

type AuthHandlerFunc func(context.Context) (context.Context, error)

type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

func authHandler(ctx context.Context, info *grpc.UnaryServerInfo, f AuthHandlerFunc) (context.Context, error) {
	if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
		return overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
	}
	return f(ctx)
}
