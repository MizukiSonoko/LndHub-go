package controller

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

type lndHubPrivateServiceServer struct{}

func (lndHubPrivateServiceServer) GetInfo(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	panic("implement me")
}

func (*lndHubPrivateServiceServer) Authorize(ctx context.Context, fullMethodName string) error {
	return nil
}

func GetLndHubPrivateServiceServer() api.LndHubPrivateServiceServer {
	return &lndHubPrivateServiceServer{}
}
