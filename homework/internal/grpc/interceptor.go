package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewLoggerServerUnaryInterceptor(l *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		l.Infow("request", "method", info.FullMethod)
		resp, err = handler(ctx, req)
		if err != nil {
			l.Infow("request not processed", "error", err)
		}

		return resp, err
	}
}

func NewLoggerServerStreamInterceptor(l *zap.SugaredLogger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		l.Infow("request", "method", info.FullMethod)
		err := handler(srv, ss)
		if err != nil {
			l.Infow("request not processed", "error", err)
		}

		return err
	}
}

func NewLoggerClientUnaryInterceptor(l *zap.SugaredLogger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		l.Infow("request", "method", method)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			l.Infow("request not processed", "error", err)
		}

		return err
	}
}

func NewLoggerClientStreamInterceptor(l *zap.SugaredLogger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		l.Infow("request", "method", method)
		stream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			l.Infow("request not processed", "error", err)
		}

		return stream, err
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
