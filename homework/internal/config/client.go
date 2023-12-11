package config

import (
	"slices"
	"time"

	"google.golang.org/grpc"
)

// DialConfig contains parameters for grpc.DialOption-s.
type DialConfig struct {
	// Blocking defines grpc.WithBlock option.
	Blocking bool `yaml:"blocking"`

	// WriteBufferSize defines grpc.WithWriteBufferSize option.
	WriteBufferSize *int `yaml:"write_buffer_size,omitempty"`

	// UserAgent defines grpc.WithUserAgent option.
	UserAgent *string `yaml:"user_agent,omitempty"`
}

// InitDialOptions parses fields in DialConfig and forms
// grpc.DialOption from them.
func (dc *DialConfig) InitDialOptions() []grpc.DialOption {
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

// LoggerConfig provides a configuration for charmbracelet/log.
type LoggerConfig struct {
	TimeFormat *string `yaml:"time_format,omitempty"`
	Output     string  `yaml:"output"`
}

// ClientConfig is the config for grpc Client.
type ClientConfig struct {
	// Addr is an address "ip:port" of the client.
	Addr               string        `yaml:"addr"`
	DialContextTimeout time.Duration `yaml:"dial_context_timeout"`
	DialConfig         DialConfig    `yaml:"dial_config"`
	Logger             LoggerConfig  `yaml:"logger"`
}
