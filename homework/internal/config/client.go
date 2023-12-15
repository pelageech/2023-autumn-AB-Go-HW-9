package config

import (
	"slices"
	"time"

	"google.golang.org/grpc"
)

// DialClient contains parameters for grpc.DialOption-s.
type DialClient struct {
	// Blocking defines grpc.WithBlock option.
	Blocking bool `yaml:"blocking"`

	// WriteBufferSize defines grpc.WithWriteBufferSize option.
	WriteBufferSize *int `yaml:"write_buffer_size,omitempty"`

	// UserAgent defines grpc.WithUserAgent option.
	UserAgent *string `yaml:"user_agent,omitempty"`
}

// InitDialOptions parses fields in DialClient and forms
// grpc.DialOption from them.
func (dc *DialClient) InitDialOptions() []grpc.DialOption {
	var s []grpc.DialOption

	if dc.Blocking {
		s = append(s, grpc.WithBlock())
	}
	if dc.UserAgent != nil {
		s = append(s, grpc.WithUserAgent(*dc.UserAgent))
	}
	if dc.WriteBufferSize != nil {
		s = append(s, grpc.WithWriteBufferSize(*dc.WriteBufferSize))
	}

	return slices.Clone(s) // copy to avoid memory leak
}

// Client is the config for grpc Client.
type Client struct {
	// Addr is an address "ip:port" of the client.
	Addr               string        `yaml:"addr"`
	DialContextTimeout time.Duration `yaml:"dial_context_timeout"`
	DialClient         DialClient    `yaml:"dial_client"`
	Logger             Logger        `yaml:"logger"`
}
