package controller

import (
	"context"
	"github.com/MizukiSonoko/LndHub-go/jwt"
	"github.com/MizukiSonoko/LndHub-go/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type lndHubServiceServer struct{}

func (*lndHubServiceServer) GetInfo(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	panic("implement me")
}

func (*lndHubServiceServer) CreateUser(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	panic("implement me")
}

func (*lndHubServiceServer) Login(ctx context.Context, req *api.LoginReq) (*api.LoginRes, error) {
	userId := req.UserId
	rawPassword := req.Password

	if userId != "mizuki" || rawPassword != "pasuwaad0" {
		return nil, status.Error(codes.Unauthenticated, "invalid")
	}
	token := jwt.GenerateToken(userId)
	log.Info("token", zap.String("token", token))
	return &api.LoginRes{
		Token: token,
	}, nil
}

func (*lndHubServiceServer) Authorize(ctx context.Context, fullMethodName string) error {
	return nil
}

func GetLndHubServiceServer() api.LndHubServiceServer {
	return &lndHubServiceServer{}
}
