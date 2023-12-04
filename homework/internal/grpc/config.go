package grpc

import (
	"slices"
	"time"

	"google.golang.org/grpc"
)

type LoggerConfig struct {
	TimeFormat *string `yaml:"time_format,omitempty"`
}

type ClientConfig struct {
	Addr               string        `yaml:"addr"`
	DialContextTimeout time.Duration `yaml:"dial_context_timeout"`
	Blocking           bool          `yaml:"blocking"`
	WriteBufferSize    *int          `yaml:"write_buffer_size,omitempty"`
	UserAgent          *string       `yaml:"user_agent,omitempty"`
	Logger             LoggerConfig  `yaml:"logger"`
}

func (cc *ClientConfig) InitDialOptions() []grpc.DialOption {
	var s []grpc.DialOption

	if cc.Blocking {
		s = append(s, grpc.WithBlock())
	}
	if cc.UserAgent != nil {
		s = append(s, grpc.WithUserAgent(*cc.UserAgent))
	}
	if cc.WriteBufferSize != nil {
		s = append(s, grpc.WithWriteBufferSize(*cc.WriteBufferSize))
	}

	return slices.Clone(s) // copy to avoid memory leak
}
