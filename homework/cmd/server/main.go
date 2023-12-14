package main

import (
	_ "embed"
	"log"
	"net"

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
		log.Fatal(err)
	}

	logger, err := cfg.Logger.Init()
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	gs := grpcinternal.NewGRPCServerPrepare(
		fileservice.New(dirfs.New(cfg.Service.Dir)),
		grpc.ChainUnaryInterceptor(
			grpcinternal.NewLoggerServerInterceptor(logger),
			grpcinternal.ValidateInterceptor,
		),
	)

	log.Printf("Server started on %s", cfg.Addr)
	if err := gs.Serve(l); err != nil {
		log.Fatalf("server set up error: %v", err)
	}
}
