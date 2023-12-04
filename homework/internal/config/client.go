package config

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpc2 "homework/internal/grpc"
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
func NewClient(ctx context.Context, config *grpc2.ClientConfig) (*Client, error) {
	logOpts := log.Options{ReportTimestamp: true}
	if config.Logger.TimeFormat != nil {
		logOpts.TimeFormat = *config.Logger.TimeFormat
	}

	loggerInterceptor := grpc.WithUnaryInterceptor(
		grpc2.NewLoggerClientInterceptor(log.NewWithOptions(os.Stdout, logOpts)),
	)
	conn, err := grpc.DialContext(ctx, config.Addr,
		append(
			config.InitDialOptions(),
			loggerInterceptor,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)...)
	if err != nil {
		return nil, fmt.Errorf("client not created: %w", err)
	}

	return &Client{conn: conn, FileServiceClient: filepb.NewFileServiceClient(conn)}, nil
}
