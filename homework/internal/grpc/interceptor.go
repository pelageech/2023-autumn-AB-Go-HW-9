package grpc

import (
	"context"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewLoggerServerInterceptor(l *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		l.Infof(`Request method: "%v"`, info.FullMethod)
		resp, err = handler(ctx, req)
		if err != nil {
			l.Infof("Error returned: %v", err)
		}

		return resp, err
	}
}

func NewLoggerClientInterceptor(l *log.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		l.Infof(`Request method: "%v"`, method)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			l.Infof("Error returned: %v", err)
		}

		return err
	}
}

// ValidateInterceptor validates incoming requests
func ValidateInterceptor(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if v, ok := req.(Validator); ok {
		if err := v.Validate(); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}

type Validator interface {
	Validate() error
}
