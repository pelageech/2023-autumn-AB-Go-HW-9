package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"homework/internal/models"
	filepb "homework/internal/proto"
	"homework/pkg/iterator"
)

const PacketSize = 4 << 10 // 4 KiB

//go:generate go run github.com/vektra/mockery/v2@v2.38.0 --name=FileService
type FileService interface {
	Ls(ctx context.Context, path models.FilePath) ([]models.FileName, error)
	Meta(ctx context.Context, path models.FilePath) (*models.FileInfo, error)
	ReadFileIterator(ctx context.Context, path models.FilePath) (iterator.Interface[[]byte], error)
}

type GRPCServer struct {
	filepb.UnimplementedFileServiceServer
	fileService FileService
}

func New(service FileService) *GRPCServer {
	return &GRPCServer{fileService: service}
}

func (s *GRPCServer) ReadFile(req *filepb.ReadFileRequest, server filepb.FileService_ReadFileServer) error {
	i, err := s.fileService.ReadFileIterator(server.Context(), req.Name)
	if err != nil {
		return status.Errorf(codes.Internal, "init iterator error: %v", err)
	}

	return iterator.Iterate(i, func(b []byte) error {
		if err := server.Send(&filepb.ReadFileReply{Stream: b}); err != nil {
			return status.Errorf(codes.Internal, "send file error: %v", err)
		}
		return nil
	})
}

func (s *GRPCServer) Ls(ctx context.Context, req *filepb.LsRequest) (*filepb.LsReply, error) {
	filenames, err := s.fileService.Ls(ctx, req.Dir)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	return &filepb.LsReply{Files: filenames}, nil
}
func (s *GRPCServer) Meta(ctx context.Context, req *filepb.MetaRequest) (*filepb.MetaReply, error) {
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
