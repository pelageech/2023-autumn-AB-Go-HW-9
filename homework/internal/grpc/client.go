package grpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"homework/internal/config"
	filepb "homework/internal/proto"
)

type Client struct {
	filepb.FileServiceClient
	conn *grpc.ClientConn
}

func (c *Client) CloseConn() error {
	return c.conn.Close()
}

// NewClient creates a new FileServiceClient from Client.
// DialContextTimeout is not used inside. The user should use it in context
// inside the function to cancel the context by themselves.
func NewClient(ctx context.Context, logger *zap.SugaredLogger, config *config.Client, otherOps ...grpc.DialOption) (*Client, error) {
	loggerInterceptor := grpc.WithUnaryInterceptor(
		NewLoggerClientInterceptor(logger),
	)

	ops := append(
		config.DialClient.InitDialOptions(),
		loggerInterceptor, // most outer interceptor
	)
	conn, err := grpc.DialContext(ctx, config.Addr, append(ops, otherOps...)...)
	if err != nil {
		return nil, fmt.Errorf("client not created: %w", err)
	}

	return &Client{conn: conn, FileServiceClient: filepb.NewFileServiceClient(conn)}, nil
}
