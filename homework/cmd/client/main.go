package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"io"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"

	"homework/internal/config"
	grpcinternal "homework/internal/grpc"
	"homework/internal/proto"
)

//go:embed config.yaml
var byteConfig []byte

func main() {
	// load config
	cfg := new(config.Client)
	if err := yaml.Unmarshal(byteConfig, cfg); err != nil {
		log.Fatal(err)
	}

	// configure logger
	logger, err := cfg.Logger.Init()
	if err != nil {
		zap.S().Fatal(err)
	}

	// connection to the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cli, err := grpcinternal.NewClient(ctx, logger, cfg,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	// take request LS
	ans, err := cli.Ls(ctx, &proto.LsRequest{Dir: "internal"})
	if err != nil {
		log.Fatal("Ls request error:", err)
	}
	log.Printf("%v\n", ans)

	// take ReadFile request
	c, _ := cli.ReadFile(ctx, &proto.ReadFileRequest{Name: "internal/fileservice/fs.go"})

	buf := new(bytes.Buffer)
	for {
		r, err := c.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("%v\n", err)
			break
		}

		buf.Write(r.GetStream())
	}
	log.Print(buf.String())
}
