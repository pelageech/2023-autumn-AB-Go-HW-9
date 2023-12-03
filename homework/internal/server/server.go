package server

import (
	"context"
	"errors"
	"io"
	"io/fs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"homework/internal/models"
	filepb "homework/internal/proto"
	hwslices "homework/pkg/slices"
)

const PacketSize = 4 << 10 // 4 KiB

type GRPCServer struct {
	filepb.UnimplementedFileServiceServer
	fileService Service
}

type Service interface {
	ReadFile(ctx context.Context, path models.FilePath) ([]byte, error)
	Ls(ctx context.Context, path models.FilePath) ([]models.FileName, error)
	Meta(ctx context.Context, path models.FilePath) (fs.FileInfo, error)
}

func (s *GRPCServer) ReadFile(req *filepb.ReadFileRequest, server filepb.FileService_ReadFileServer) error {
	if req.Name == "" {
		return status.Errorf(codes.InvalidArgument, "file path is empty")
	}

	b, err := s.fileService.ReadFile(server.Context(), req.Name)
	if err != nil && !errors.Is(err, io.EOF) {
		return status.Errorf(codes.Internal, "read file error: %v", err)
	}

	packets := hwslices.Split(b, PacketSize)

	for _, p := range packets {
		if err := server.Send(&filepb.ReadFileReply{Stream: p}); err != nil {
			return status.Errorf(codes.Internal, "send file error: %v", err)
		}
	}

	return nil
}

func (s *GRPCServer) Ls(ctx context.Context, req *filepb.LsRequest) (*filepb.LsReply, error) {
	if req.Dir == "" {
		return nil, status.Errorf(codes.InvalidArgument, "file path is empty")
	}

	filenames, err := s.fileService.Ls(ctx, req.Dir)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	return &filepb.LsReply{Files: filenames}, nil
}
func (s *GRPCServer) Meta(ctx context.Context, req *filepb.MetaRequest) (*filepb.MetaReply, error) {
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "file path is empty")
	}

	stat, err := s.fileService.Meta(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err)
	}

	metaReply := &filepb.MetaReply{
		Size:  stat.Size(),
		Mode:  uint32(stat.Mode()),
		IsDir: stat.IsDir(),
	}

	return metaReply, nil
}
