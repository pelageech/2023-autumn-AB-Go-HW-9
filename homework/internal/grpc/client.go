package grpc

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
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

// NewClient creates a new FileServiceClient from ClientConfig.
// DialContextTimeout is not used inside. The user should use it in context
// inside the function to cancel the context by themselves.
func NewClient(ctx context.Context, logger *log.Logger, config *config.ClientConfig) (*Client, error) {
	loggerInterceptor := grpc.WithUnaryInterceptor(
		NewLoggerClientInterceptor(logger),
	)

	conn, err := grpc.DialContext(ctx, config.Addr,
		append(
			config.DialConfig.InitDialOptions(),
			loggerInterceptor,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)...)
	if err != nil {
		return nil, fmt.Errorf("client not created: %w", err)
	}

	return &Client{conn: conn, FileServiceClient: filepb.NewFileServiceClient(conn)}, nil
}
