package controller

import (
	"context"
	"encoding/json"
	"github.com/MizukiSonoko/LndHub-go/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"net/http"
	"strconv"
)

type lndHubServiceServer struct{}

func (*lndHubServiceServer) GetInfo(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	panic("implement me")
}

func (*lndHubServiceServer) CreateUser(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	panic("implement me")
}

func (*lndHubServiceServer) Login(ctx context.Context, req *api.LoginReq) (*api.LoginRes, error) {
	pUserId := r.Form.Get("userId")
	pPassword := r.Form.Get("password")

	amount, err := strconv.Atoi(pAmount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "amount should be number"})
		return
	}
	if amount < 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "amount should be plus"})
		return
	}

	token := middleware.GenerateToken(nil)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		TokenResp{Token: token})
}

func (*lndHubServiceServer) Authorize(ctx context.Context, fullMethodName string) error {
	return nil
}

func GetLndHubServiceServer() api.LndHubServiceServer {
	return &lndHubServiceServer{}
}
