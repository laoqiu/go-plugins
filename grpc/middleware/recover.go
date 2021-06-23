package middleware

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecoveryHandlerFunc func(context.Context, interface{}) error

func recoverFrom(ctx context.Context, p interface{}, f RecoveryHandlerFunc) error {
	if f == nil {
		return status.Errorf(codes.Internal, "%v", p)
	}
	return f(ctx, p)
}
