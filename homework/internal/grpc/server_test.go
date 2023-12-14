package grpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"homework/internal/config"

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

	client *Client
}

func (s *grpcSuite) SetupSuite() {
	// config server
	s.fileService = &mocks.FileService{}

	s.listener = bufconn.Listen(1 << 20)
	s.T().Log("listener configured")

	s.grpcserver = NewGRPCServerPrepare(s.fileService)
	s.T().Log("server configured")

	// config client
	cfg := &config.Client{Addr: "bufnet"}
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return s.listener.Dial()
	}
	client, err := NewClient(context.Background(), zap.S(), cfg,
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.Error(err)
	}
	s.T().Log("client configured")

	s.client = client
	s.T().Log("client configured")

	// setup server
	go func() {
		err = s.grpcserver.Serve(s.listener)
		if err != nil {
			s.Error(err)
		}
	}()
}

func (s *grpcSuite) TearDownSuite() {
	err := s.client.CloseConn()
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

const maxReadIterationsCount = 3

type badReader struct {
	r io.Reader
	i int
}

func (r *badReader) Read(b []byte) (int, error) {
	if r.i >= maxReadIterationsCount {
		return 0, fmt.Errorf("test error")
	}
	r.i++

	return r.r.Read(b)
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
		{"error while reading", args{name: "bad.wncry"}, nil, true},
		{"empty file", args{name: "empty.txt"}, []byte{}, false},
		{"small file", args{name: "small.txt"}, b[:1<<10], false},
		{"ordinary file", args{name: "ordinary.txt"}, b[:8<<10], false},
		{"large file", args{name: "large.txt"}, b, false},
	}

	s.fileService.EXPECT().ReadFileIterator(mock.Anything, tests[0].args.name).
		Return(nil, fs.ErrInvalid).Once()

	s.fileService.EXPECT().ReadFileIterator(mock.Anything, tests[1].args.name).
		Return(iterator.NewReaderIterator(context.Background(), io.NopCloser(&badReader{r: bytes.NewReader(b[:6])}), make([]byte, 1)), nil).Once()

	for _, t := range tests[2:] {
		s.fileService.EXPECT().ReadFileIterator(mock.Anything, t.args.name).
			Return(iterator.NewReaderIterator(context.Background(), io.NopCloser(bytes.NewReader(t.want)), make([]byte, 4<<10)), nil).Once()
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt := tt
			s.T().Parallel()

			var r *filepb.ReadFileReply
			buf := bytes.NewBuffer([]byte{})

			cli, err := s.client.ReadFile(context.Background(), &filepb.ReadFileRequest{Name: tt.args.name})
			if err == nil {
				for {
					r, err = cli.Recv()
					if errors.Is(err, io.EOF) {
						err = nil
						break
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

	s.fileService.EXPECT().Ls(mock.Anything, tests[0].args.name).
		Return(nil, fs.ErrInvalid).Once()
	for _, t := range tests[1:] {
		s.fileService.EXPECT().Ls(mock.Anything, t.args.name).
			Return(t.want, nil).Once()
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt := tt
			s.T().Parallel()

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

	s.fileService.EXPECT().Meta(mock.Anything, tests[0].args.name).
		Return(nil, fs.ErrInvalid).Once()
	for _, t := range tests[1:] {
		s.fileService.EXPECT().Meta(mock.Anything, t.args.name).
			Return(t.want, nil).Once()
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt := tt
			s.T().Parallel()

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
