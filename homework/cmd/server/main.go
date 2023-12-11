package main

import (
	"github.com/charmbracelet/log"
	"google.golang.org/grpc"

	"homework/internal/fileservice"
	grpcinternal "homework/internal/grpc"
	"homework/internal/repo/dirfs"
)

func main() {
	gs, l, err := grpcinternal.NewGRPCServerPrepare(":50051",
		fileservice.New(dirfs.New(".")),
		grpc.ChainUnaryInterceptor(
			grpcinternal.NewLoggerServerInterceptor(log.Default()),
			grpcinternal.ValidateInterceptor,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server started at :50051\n")
	if err := gs.Serve(l); err != nil {
		log.Fatal(err)
	}
}
