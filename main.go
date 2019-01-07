package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ToDo: make constructor func and hide member because we should not change it.
type IndInfo struct {
	LightningId     string
	IdentityAddress string
	IdentityPubkey  string
}

type Lnd interface {
	GetInfo() (IndInfo, error)
}

type lndClient struct {
	client lnrpc.LightningClient
	c      context.Context
}

func (l *lndClient) GetInfo() (IndInfo, error) {
	res, err := l.client.GetInfo(l.c, &lnrpc.GetInfoRequest{})
	if err != nil {
		return IndInfo{}, fmt.Errorf("grpc error err:%s", err)
	}
	return IndInfo{
		res.LightningId,
		res.IdentityAddress,
		res.IdentityPubkey,
	}, nil
}

func getClient(address, tlsPath string) *grpc.ClientConn {
	_ = os.Setenv("GRPC_SSL_CIPHER_SUITES", "HIGH+ECDSA")

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

func newLnd(address, certPath string) Lnd {
	return &lndClient{
		lnrpc.NewLightningClient(getClient(address, certPath)),
		context.Background(),
	}
}

func main() {
	// I use macOS
	lnd := newLnd("localhost:10009", os.Getenv("HOME")+"/Library/Application Support/Lnd/tls.cert")
	{
		res, err := lnd.GetInfo()
		if err != nil {
			fmt.Printf("GetInfo failed err:%s", err)
			os.Exit(1)
		}
		fmt.Printf("lightningId:%s\n", res.LightningId)
		fmt.Printf("IdentityAddress:%s\n", res.IdentityAddress)
		fmt.Printf("IdentityPubkey:%s\n", res.IdentityPubkey)
	}
}
