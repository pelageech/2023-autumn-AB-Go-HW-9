package main

import (
	_ "embed"
	"flag"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"homework/internal/config"
	"homework/internal/fileservice"
	grpcinternal "homework/internal/grpc"
	"homework/internal/repo/dirfs"
)

var configPath = flag.String("config-path", "config.yaml", "Client file configuration")

func main() {
	flag.Parse()

	f, err := os.Open(*configPath)
	if err != nil {
		log.Fatalln("file open error:", err)
	}
	defer func() {
		_ = f.Close()
	}()

	byteConfig, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln("file read error:", err)
	}

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
			grpcinternal.NewLoggerServerUnaryInterceptor(logger),
			grpcinternal.ValidateInterceptor,
		),
		grpc.StreamInterceptor(grpcinternal.NewLoggerServerStreamInterceptor(logger)),
	)

	log.Printf("Server started on %s", cfg.Addr)
	if err := gs.Serve(l); err != nil {
		log.Fatalf("server set up error: %v", err)
	}
}
