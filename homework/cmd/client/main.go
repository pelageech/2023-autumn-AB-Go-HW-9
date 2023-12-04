package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"

	"homework/internal/config"
	grpcinternal "homework/internal/grpc"
	"homework/internal/proto"
)

//go:embed config.yaml
var byteConfig []byte

func main() {
	cfg := new(config.ClientConfig)
	if err := yaml.Unmarshal(byteConfig, cfg); err != nil {
		log.Fatal(err)
	}

	logOpts := log.Options{ReportTimestamp: true}
	if cfg.Logger.TimeFormat != nil {
		logOpts.TimeFormat = *cfg.Logger.TimeFormat
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cli, err := grpcinternal.NewClient(ctx, log.NewWithOptions(os.Stdout, logOpts), cfg)
	if err != nil {
		log.Fatal(err)
	}

	ans, err := cli.Ls(ctx, &proto.LsRequest{Dir: "internal"})
	log.Printf("%v\n", ans)
	c, err := cli.ReadFile(ctx, &proto.ReadFileRequest{Name: "internal/fileservice/fs.go"})

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
