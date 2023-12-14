package grpc

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

type ClientOp func(*Client)

// NewClient creates a new FileServiceClient from Client.
// DialContextTimeout is not used inside. The user should use it in context
// inside the function to cancel the context by themselves.
func NewClient(ctx context.Context, logger *zap.SugaredLogger, config *config.Client) (*Client, error) {
	loggerInterceptor := grpc.WithUnaryInterceptor(
		NewLoggerClientInterceptor(logger),
	)

	conn, err := grpc.DialContext(ctx, config.Addr,
		append(
			config.DialClient.InitDialOptions(),
			loggerInterceptor,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)...)
	if err != nil {
		return nil, fmt.Errorf("client not created: %w", err)
	}

	return &Client{conn: conn, FileServiceClient: filepb.NewFileServiceClient(conn)}, nil
}
