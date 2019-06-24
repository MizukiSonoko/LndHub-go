package main

import (
	"context"
	"fmt"
	"github.com/MizukiSonoko/LndHub-go/protobuf"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	target := "localhost:50000"
	fmt.Printf("target:\"%s\"\n", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("client connection error:%s\n", err)
	}
	defer conn.Close()
	ctx := context.Background()

	{
		client := api.NewLndHubServiceClient(conn)
		res, err := client.Login(ctx, &api.LoginReq{
			UserId:   "mizuki",
			Password: "pasuwaad0",
		})
		if err != nil {
			fmt.Printf("Login failed err:%s\n", err)
			return
		}
		fmt.Printf("login successful! token:%s\n", res.Token)
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", res.Token)
	}
	{
		client := api.NewLndHubPrivateServiceClient(conn)
		res, err := client.GetBtc(ctx, &empty.Empty{})
		if err != nil {
			fmt.Printf("Login failed err:%s\n", err)
			return
		}
		fmt.Printf("BTC address:%s\n", res.Address)
	}
}
