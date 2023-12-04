package main

import (
	"io/fs"
	"os"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc"

	"homework/internal/fileservice"
	grpcinternal "homework/internal/grpc"
)

func main() {
	gs, l, err := grpcinternal.NewGRPCServerPrepare(":50051",
		fileservice.New(os.DirFS(".").(fs.ReadDirFS)),
		grpc.ChainUnaryInterceptor(
			grpcinternal.NewLoggerServerInterceptor(log.Default()),
			grpcinternal.ValidateInterceptor,
		),
		grpc.StreamInterceptor(grpcinternal.StreamValidateInterceptor),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server started at :50051\n")
	if err := gs.Serve(l); err != nil {
		log.Fatal(err)
	}
}
