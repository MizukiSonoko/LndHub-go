package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func getClient(address, tlsPath string) *grpc.ClientConn {

	creds, err := credentials.NewClientTLSFromFile(tlsPath, "")
	if err != nil {
		fmt.Printf("NewClientTLSFromFile failed %v", err)
		os.Exit(1)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		fmt.Printf("unable to connect to RPC server: %v", err)
		os.Exit(1)
	}

	return conn
}

func main() {
	_ = os.Setenv("GRPC_SSL_CIPHER_SUITES", "HIGH+ECDSA")

	// I use macOS
	client := lnrpc.NewLightningClient(getClient(
		"localhost:10009", os.Getenv("HOME")+"/Library/Application Support/Lnd/tls.cert"))
	res, err := client.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	if err != nil {
		fmt.Printf("GetInfo failed err:%s", err)
		os.Exit(1)
	}
	fmt.Printf("lightningId:%s\n", res.LightningId)
}
