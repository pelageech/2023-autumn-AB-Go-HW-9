package main

import (
	_ "embed"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"homework/internal/config"
	"homework/internal/fileservice"
	grpcinternal "homework/internal/grpc"
	"homework/internal/repo/dirfs"
)

//go:embed config.yaml
var byteConfig []byte

func main() {
	cfg := new(config.Server)
	if err := yaml.Unmarshal(byteConfig, cfg); err != nil {
		zap.S().Fatal(err)
	}

	logger, err := cfg.Logger.Init()
	if err != nil {
		zap.S().Fatal(err)
	}

	gs, l, err := grpcinternal.NewGRPCServerPrepare(cfg.Addr,
		fileservice.New(dirfs.New(cfg.Service.Dir)),
		grpc.ChainUnaryInterceptor(
			grpcinternal.NewLoggerServerInterceptor(logger),
			grpcinternal.ValidateInterceptor,
		),
	)
	if err != nil {
		log.Fatalf("server prepare error: %v", err)
	}

	log.Printf("Server started on %s", cfg.Addr)
	if err := gs.Serve(l); err != nil {
		log.Fatalf("server set up error: %v", err)
	}
}
