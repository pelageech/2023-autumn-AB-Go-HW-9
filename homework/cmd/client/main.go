package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"

	"homework/internal/config"
	grpcinternal "homework/internal/grpc"
	"homework/internal/proto"
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
