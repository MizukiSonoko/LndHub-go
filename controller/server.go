package controller

import (
	"context"
	"fmt"
	"github.com/MizukiSonoko/lnd-gateway/interceptor"
	"github.com/MizukiSonoko/lnd-gateway/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net"
	"time"
)

var (
	log           = logger.NewLogger()
	ServerIsReady = make(chan struct{})
)

type GRPCServer struct {
	*grpc.Server
}

func NewGRPCServer(ctx context.Context) *GRPCServer {
	opts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}
	zapLogger := logger.NewLogger()

	grpc_zap.ReplaceGrpcLogger(zapLogger)

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(zapLogger, opts...),
			interceptor.UnaryAuthenticateInterceptor(),
			interceptor.UnaryAuthorizationInterceptor(),
		),
	)
	return &GRPCServer{s}
}

func (s *GRPCServer) Start() error {
	defer close(ServerIsReady)

	// ToDo use config
	address := fmt.Sprintf("%s:%d",
		"0.0.0.0", 50000)

	log.Info("rpc server listen", zap.String("address", address))
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrapf(err, "listen failed", "tcp", address)
	}

	return s.Serve(listen)
}

func (s *GRPCServer) Stop() error {
	s.GracefulStop()
	return nil
}
