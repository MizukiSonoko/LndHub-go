package controller

import (
	"context"
	"github.com/MizukiSonoko/lnd-gateway/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
)

type lndHubServiceServer struct{}

func (*lndHubServiceServer) GetInfo(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	panic("implement me")
}

func (*lndHubServiceServer) CreateUser(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	panic("implement me")
}

func (*lndHubServiceServer) Login(ctx context.Context, req *api.LoginReq) (*api.LoginRes, error) {
	panic("implement me")
}

func (*lndHubServiceServer) Authorize(ctx context.Context, fullMethodName string) error {
	return nil
}

func GetLndHubServiceServer() api.LndHubServiceServer {
	return &lndHubServiceServer{}
}
