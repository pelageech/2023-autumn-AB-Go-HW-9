package grpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"homework/internal/grpc/mocks"
	"homework/internal/models"
	filepb "homework/internal/proto"
	"homework/pkg/iterator"
)

type grpcSuite struct {
	suite.Suite
	fileService *mocks.FileService
	grpcserver  *grpc.Server
	listener    *bufconn.Listener
	conn        *grpc.ClientConn

	server filepb.FileServiceServer
	client filepb.FileServiceClient
}

func (s *grpcSuite) SetupSuite() {
	s.fileService = &mocks.FileService{}
	s.server = NewFileServiceServer(s.fileService)
	s.listener = bufconn.Listen(1 << 20)
	s.T().Log("listener configured")
	s.grpcserver = grpc.NewServer()
	s.T().Log("server configured")

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return s.listener.Dial()
	}
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.Error(err)
	}
	s.T().Log("conn configured")

	s.conn = conn
	s.client = filepb.NewFileServiceClient(conn)
	s.T().Log("client configured")

	filepb.RegisterFileServiceServer(s.grpcserver, s.server)
	s.T().Log("server registered")

	go func() {
		err = s.grpcserver.Serve(s.listener)
		if err != nil {
			s.Error(err)
		}
	}()
}

func (s *grpcSuite) TearDownSuite() {
	err := s.conn.Close()
	if err != nil {
		s.Error(err)
	}
	s.T().Log("conn closed")

	s.grpcserver.GracefulStop()
	s.T().Log("server closed")

	err = s.listener.Close()
	if err != nil {
		s.Error(err)
	}
	s.T().Log("listener closed")
}

func (s *grpcSuite) TestGRPCServer_ReadFileIterator() {
	rndm := rand.New(rand.NewSource(42))
	b := make([]byte, 8<<20)
	_, _ = rndm.Read(b)
	type args struct {
		name models.FilePath
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"invalid path", args{name: ":.<%fkv'`"}, nil, true},
		{"empty file", args{name: "empty.txt"}, []byte{}, false},
		{"small file", args{name: "small.txt"}, b[:1<<10], false},
		{"ordinary file", args{name: "ordinary.txt"}, b[:8<<10], false},
		{"large file", args{name: "large.txt"}, b, false},
	}

	s.fileService.On("ReadFileIterator", mock.Anything, tests[0].args.name).
		Return(nil, fs.ErrInvalid).Once()
	for _, t := range tests[1:] {
		s.fileService.On("ReadFileIterator", mock.Anything, t.args.name).
			Return(iterator.NewReaderIterator(context.Background(), io.NopCloser(bytes.NewReader(t.want)), 4<<10), nil).Once()
	}

	var r *filepb.ReadFileReply
	buf := bytes.NewBuffer([]byte{})

	for _, tt := range tests {
		buf.Reset()
		s.Run(tt.name, func() {
			cli, err := s.client.ReadFile(context.Background(), &filepb.ReadFileRequest{Name: tt.args.name})
			if err == nil {
				for {
					r, err = cli.Recv()
					if errors.Is(err, io.EOF) {
						err = nil
						return
					}
					if err != nil {
						break
					}
					buf.Write(r.GetStream())
				}
			}

			if (err != nil) != tt.wantErr {
				s.T().Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got := buf.Bytes()

			assert.Equal(s.T(), tt.want, got)
		})
	}

}

func (s *grpcSuite) TestGRPCServer_Ls() {
	type args struct {
		name models.FilePath
	}
	tests := []struct {
		name    string
		args    args
		want    []models.FileName
		wantErr bool
	}{
		{"invalid path", args{name: ":.<%fkv'`"}, nil, true},
		{"empty dir", args{name: "empty"}, []models.FileName{}, false},
		{"one file", args{name: "one"}, []models.FileName{"aboba"}, false},
		{"many files", args{name: "many"}, []models.FileName{"aboba1", "aboba2", "aboba3"}, false},
	}

	s.fileService.On("Ls", mock.Anything, tests[0].args.name).
		Return(nil, fs.ErrInvalid).Once()
	for _, t := range tests[1:] {
		s.fileService.On("Ls", mock.Anything, t.args.name).
			Return(t.want, nil).Once()
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			cli, err := s.client.Ls(context.Background(), &filepb.LsRequest{Dir: tt.args.name})

			if (err != nil) != tt.wantErr {
				s.T().Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got := cli.GetFiles()

			assert.ElementsMatch(s.T(), tt.want, got)
		})
	}
}

func (s *grpcSuite) TestGRPCServer_Meta() {
	type args struct {
		name models.FilePath
	}
	tests := []struct {
		name    string
		args    args
		want    *models.FileInfo
		wantErr bool
	}{
		{"invalid path", args{name: ""}, nil, true},
		{"dir", args{name: "empty"}, &models.FileInfo{Size: 0, Mode: 0o777, IsDir: true}, false},
		{"one file", args{name: "one"}, &models.FileInfo{Size: 1987, Mode: 0o145, IsDir: false}, false},
	}

	s.fileService.On("Meta", mock.Anything, tests[0].args.name).
		Return(nil, fs.ErrInvalid).Once()
	for _, t := range tests[1:] {
		s.fileService.On("Meta", mock.Anything, t.args.name).
			Return(t.want, nil).Once()
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			r, err := s.client.Meta(context.Background(), &filepb.MetaRequest{Name: tt.args.name})

			if (err != nil) != tt.wantErr {
				s.T().Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			assert.Equal(s.T(), tt.want.IsDir, r.IsDir)
			assert.Equal(s.T(), tt.want.Size, r.Size)
			assert.Equal(s.T(), tt.want.Mode, r.Mode)
		})
	}
}

func TestGRPCServer_Suite(t *testing.T) {
	suite.Run(t, new(grpcSuite))
}