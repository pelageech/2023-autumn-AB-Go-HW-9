package grpc

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"homework/internal/config"
	"homework/internal/grpc/mocks"
	"homework/internal/proto"
	"homework/pkg/iterator"
)

type grpcInterceptorSuite struct {
	suite.Suite
	fileService *mocks.FileService
	grpcserver  *grpc.Server
	listener    *bufconn.Listener

	client *Client

	logsClient *observer.ObservedLogs
	logsServer *observer.ObservedLogs
}

func (s *grpcInterceptorSuite) SetupSuite() {
	// server config
	s.fileService = mocks.NewFileService(s.T())

	s.listener = bufconn.Listen(1 << 20)
	s.T().Log("listener configured")

	coreServer, logsServer := observer.New(zap.InfoLevel)
	s.logsServer = logsServer
	loggerServer := zap.New(coreServer).Sugar()

	s.grpcserver = NewGRPCServerPrepare(s.fileService,
		grpc.UnaryInterceptor(NewLoggerServerUnaryInterceptor(loggerServer)),
		grpc.StreamInterceptor(NewLoggerServerStreamInterceptor(loggerServer)),
	)
	s.T().Log("server configured")

	// client config
	coreClient, logsClient := observer.New(zap.InfoLevel)
	s.logsClient = logsClient
	loggerClient := zap.New(coreClient).Sugar()

	cfg := &config.Client{Addr: "bufnet"}
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return s.listener.Dial()
	}
	client, err := NewClient(context.Background(), loggerClient, cfg,
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.Error(err)
	}

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

func (s *grpcInterceptorSuite) TearDownSuite() {
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

func (s *grpcInterceptorSuite) TestLoggerInterceptor_ReadFileSuccess() {
	s.fileService.EXPECT().ReadFileIterator(mock.Anything, "good").
		Return(iterator.NewReaderIterator(context.Background(), io.NopCloser(bytes.NewReader([]byte{54, 27, 0})), []byte{0}), nil).
		Once()
	cli, err := s.client.ReadFile(context.Background(), &proto.ReadFileRequest{Name: "good"})
	assert.NoError(s.T(), err)

	buf := new(bytes.Buffer)
	for {
		repl, err := cli.Recv()
		if err == io.EOF {
			break
		}
		assert.NoError(s.T(), err)
		buf.Write(repl.GetStream())
	}
	assert.Equal(s.T(), []byte{54, 27, 0}, buf.Bytes())

	message := "request"
	method := "/file.FileService/ReadFile"

	logsServer := s.logsServer.TakeAll()[0]
	assert.Equal(s.T(), message, logsServer.Message)
	actualMethod := logsServer.Context[0].String
	assert.Equal(s.T(), method, actualMethod)

	logsClient := s.logsClient.TakeAll()[0]
	assert.Equal(s.T(), message, logsClient.Message)
	actualMethod = logsServer.Context[0].String
	assert.Equal(s.T(), method, actualMethod)
}

func (s *grpcInterceptorSuite) TestLoggerInterceptor_LsSuccess() {
	expectedLs := []string{"a", "b", "c"}

	s.fileService.EXPECT().Ls(mock.Anything, "good").
		Return(expectedLs, nil).
		Once()

	repl, err := s.client.Ls(context.Background(), &proto.LsRequest{Dir: "good"})
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), expectedLs, repl.GetFiles())

	message := "request"
	method := "/file.FileService/Ls"

	logsServer := s.logsServer.TakeAll()[0]
	assert.Equal(s.T(), message, logsServer.Message)
	actualMethod := logsServer.Context[0].String
	assert.Equal(s.T(), method, actualMethod)

	logsClient := s.logsClient.TakeAll()[0]
	assert.Equal(s.T(), message, logsClient.Message)
	actualMethod = logsServer.Context[0].String
	assert.Equal(s.T(), method, actualMethod)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(grpcInterceptorSuite))
}
