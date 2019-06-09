package lightning

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
)

type IndInfo struct {
	LightningId     string
	IdentityAddress string
	IdentityPubkey  string
}

type AddInvoiceResp struct {
	PaymentRequest string
}

type DecodePayResp struct {
	NumSatoshis uint
	Description string
	PaymentHash string
}

type SendRequest struct {
}

type SendResponse struct {
	PaymentError    string
	PaymentPreimage []byte
}

type Lnd interface {
	GetInfo() (IndInfo, error)
	NewAddress() (string, error)
	AddInvoice(memo string, amt uint) (AddInvoiceResp, error)
	DecodePay(invoice string) (DecodePayResp, error)
	ConnectPeer(publicKey, host string) error
	GetSendPaymentClient() (SendPaymentClient, error)
}

type SendPaymentClient struct {
	client lnrpc.Lightning_SendPaymentClient
}

func (c *SendPaymentClient) Receive(f func(resp SendResponse, err error)) {
	resp, err := c.client.Recv()

	f(SendResponse{
		PaymentError:    resp.PaymentError,
		PaymentPreimage: resp.PaymentPreimage,
	}, err)
}

func (c *SendPaymentClient) Send(req SendRequest) error {
	return c.client.Send(&lnrpc.SendRequest{})
}

type lndClient struct {
	client lnrpc.LightningClient
	c      context.Context
}

func (l *lndClient) NewAddress() (string, error) {
	res, err := l.client.NewAddress(l.c, &lnrpc.NewAddressRequest{
		// Ref: https://github.com/BlueWallet/LndHub/blob/master//class/User.js#L107
		Type: lnrpc.AddressType_WITNESS_PUBKEY_HASH,
	})
	if err != nil {
		return "", fmt.Errorf("NewAddress failed err:%s", err)
	}
	return res.Address, nil
}

func (l *lndClient) GetInfo() (IndInfo, error) {
	res, err := l.client.GetInfo(l.c, &lnrpc.GetInfoRequest{})
	if err != nil {
		return IndInfo{}, fmt.Errorf("grpc error err:%s", err)
	}
	return IndInfo{
		LightningId:    res.Alias,
		IdentityPubkey: res.IdentityPubkey,
	}, nil
}

func (l *lndClient) ConnectPeer(publicKey, host string) error {
	_, err := l.client.ConnectPeer(l.c, &lnrpc.ConnectPeerRequest{
		Addr: &lnrpc.LightningAddress{
			Pubkey: publicKey,
			Host:   host,
		},
	})
	if err != nil {
		return fmt.Errorf("grpc error err:%s", err)
	}
	return nil
}

func (l *lndClient) AddInvoice(memo string, amt uint) (AddInvoiceResp, error) {
	res, err := l.client.AddInvoice(l.c, &lnrpc.Invoice{
		Memo:  memo,
		Value: int64(amt),
	})
	if err != nil {
		return AddInvoiceResp{}, fmt.Errorf("grpc error err:%s", err)
	}
	return AddInvoiceResp{
		PaymentRequest: res.PaymentRequest,
	}, nil
}

func (l *lndClient) DecodePay(invoice string) (DecodePayResp, error) {
	res, err := l.client.DecodePayReq(l.c, &lnrpc.PayReqString{
		PayReq: invoice,
	})
	if err != nil {
		return DecodePayResp{}, fmt.Errorf("grpc error err:%s", err)
	}

	return DecodePayResp{
		NumSatoshis: uint(res.NumSatoshis),
		Description: res.Description,
		PaymentHash: res.PaymentHash,
	}, nil
}

func (l *lndClient) GetSendPaymentClient() (SendPaymentClient, error) {
	res, err := l.client.SendPayment(l.c)
	if err != nil {
		return SendPaymentClient{}, fmt.Errorf("grpc error err:%s", err)
	}
	return SendPaymentClient{client: res}, nil
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

func NewLnd(address, certPath string) Lnd {
	return &lndClient{
		lnrpc.NewLightningClient(getClient(address, certPath)),
		context.Background(),
	}
}
