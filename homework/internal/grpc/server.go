package grpc

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"homework/internal/models"
	filepb "homework/internal/proto"
	"homework/pkg/iterator"
)

const PacketSize = 4 << 10 // 4 KiB

// FileService provides some functions for working with fs.File
// items.
// All the methods support context.
//
//go:generate go run github.com/vektra/mockery/v2@v2.38.0 --name=FileService
type FileService interface {
	// Ls returns a list of files containing in the given path.
	Ls(ctx context.Context, path models.FilePath) ([]models.FileName, error)

	// Meta returns meta-data of the file in path if exists.
	Meta(ctx context.Context, path models.FilePath) (*models.FileInfo, error)

	// ReadFileIterator opens a file on FS and returns an iterator *iterator.ReaderIterator.
	ReadFileIterator(ctx context.Context, path models.FilePath) (*iterator.ReaderIterator, error)
}

type Server struct {
	filepb.UnimplementedFileServiceServer
	fileService FileService
}

func NewFileServiceServer(service FileService) *Server {
	return &Server{fileService: service}
}

func NewGRPCServerPrepare(service FileService, op ...grpc.ServerOption) *grpc.Server {
	gs := grpc.NewServer(op...)
	filepb.RegisterFileServiceServer(gs, NewFileServiceServer(service))

	return gs
}

func (s *Server) ReadFile(req *filepb.ReadFileRequest, server filepb.FileService_ReadFileServer) error {
	i, err := s.fileService.ReadFileIterator(server.Context(), req.Name)
	if err != nil {
		return status.Errorf(codes.Internal, "init iterator error: %v", err)
	}
	defer func() {
		f := i.Reader().(io.ReadCloser)
		_ = f.Close()
	}()

	return iterator.Iterate(iterator.Simple[[]byte](i), func(b []byte) error {
		if err := server.Send(&filepb.ReadFileReply{Stream: b}); err != nil {
			return status.Errorf(codes.Internal, "send file error: %v", err)
		}
		return nil
	})
}

func (s *Server) Ls(ctx context.Context, req *filepb.LsRequest) (*filepb.LsReply, error) {
	filenames, err := s.fileService.Ls(ctx, req.Dir)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	return &filepb.LsReply{Files: filenames}, nil
}
func (s *Server) Meta(ctx context.Context, req *filepb.MetaRequest) (*filepb.MetaReply, error) {
	stat, err := s.fileService.Meta(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	metaReply := &filepb.MetaReply{
		Size:  stat.Size,
		Mode:  stat.Mode,
		IsDir: stat.IsDir,
	}

	return metaReply, nil
}
