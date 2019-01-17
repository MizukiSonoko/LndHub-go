package main

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
)

// ToDo: make constructor func and hide member because we should not change it.
type IndInfo struct {
	LightningId     string
	IdentityAddress string
	IdentityPubkey  string
}

type Lnd interface {
	GetInfo() (IndInfo, error)
	ConnectPeer(publicKey, host string) error
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

func (l *lndClient) ConnectPeer(publicKey, host string) error {
	_, err := l.client.ConnectPeer(l.c, &lnrpc.ConnectPeerRequest{
		Addr: &lnrpc.LightningAddress{
			PubKeyHash: publicKey,
			Host:       host,
		},
	})
	if err != nil {
		return fmt.Errorf("grpc error err:%s", err)
	}
	return nil
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

// Ref: https://github.com/btcsuite/btcd/blob/86fed781132ac890ee03e906e4ecd5d6fa180c64/rpcclient/examples/bitcoincorehttp/main.go
func getBlockchainInfo(){
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18555",
		User:         "sonokko",
		Pass:         "1qazxsw23edcvfr4",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		panic(err)
	}
	defer client.Shutdown()

	blockCount, err := client.GetBlockCount()
	if err != nil {
		panic(err)
	}
	log.Printf("Block count: %d", blockCount)
}

func newLnd(address, certPath string) Lnd {
	return &lndClient{
		lnrpc.NewLightningClient(getClient(address, certPath)),
		context.Background(),
	}
}

func main() {
	getBlockchainInfo()
	return
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
