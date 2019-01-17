package main

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/database"
	_ "github.com/btcsuite/btcd/database/ffldb"
	"os"
	"encoding/hex"


	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/btcsuite/btcd/blockchain"
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

func getBlockchainInfo(){

	dbPath := "/Users/mizuki/Library/Application Support/Btcd/data/simnet/blocks_ffldb"
	_ = os.RemoveAll(dbPath)
	db, err := database.Create("ffldb", dbPath, chaincfg.MainNetParams.Net)
	if err != nil {
		fmt.Printf("Failed to create database: %v\n", err)
		return
	}
	defer os.RemoveAll(dbPath)
	defer db.Close()

	chain, err := blockchain.New(&blockchain.Config{
		DB:          db,
		ChainParams: &chaincfg.MainNetParams,
		TimeSource:  blockchain.NewMedianTime(),
	})
	if err != nil{
		panic(chain)
	}

	hashStr := "683e86bd5c6d110d91b94b97137ba6bfe02dbbdb8e3dff722a669b5d69d77af6"
	hash, err := hex.DecodeString(hashStr)
	if err != nil{
		panic(err)
	}
	h := new(chainhash.Hash)
	err =h.SetBytes(hash)
	if err != nil{
		panic(err)
	}
	res, err := chain.HaveBlock(h)
	if err != nil{
		panic(err)
	}
	print(res)

	head, err := chain.BlockHashByHeight(0)
	if err != nil{
		panic(err)
	}
	fmt.Printf("head:%s\n", head.String())
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
