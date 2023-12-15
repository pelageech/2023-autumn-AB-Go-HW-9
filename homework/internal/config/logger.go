package config

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger provides a configuration for uber/zap.
type Logger struct {
	Output []string `yaml:"output"`
}

func (l *Logger) Init() (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = l.Output

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed logger build: %w", err)
	}

	return logger.Sugar(), nil
}
