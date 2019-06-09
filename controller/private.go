package controller

import (
	"context"
	"github.com/MizukiSonoko/LndHub-go/bitcoin"
	"github.com/MizukiSonoko/LndHub-go/entity"
	"github.com/MizukiSonoko/LndHub-go/lightning"
	"github.com/MizukiSonoko/LndHub-go/logger"
	"github.com/MizukiSonoko/LndHub-go/protobuf"
	"github.com/MizukiSonoko/LndHub-go/repository"
	"github.com/golang/protobuf/ptypes/empty"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"os"
	"time"
)

const (
	CtxUserIdKey = "userId"
)

var (
	repo           repository.UserRepo
	lnd            lightning.Lnd
	bc             bitcoin.BitcoinClient
	identityPubkey string

	log = logger.NewLogger()
)

func init() {
	repo = repository.NewUserRepo()
	// In now, using macOS
	lnd = lightning.NewLnd(
		"localhost:10009",
		os.Getenv("HOME")+"/.lnd/tls.cert")
	info, err := lnd.GetInfo()
	if err != nil {
		log.Fatal("lnd GetInfo failed", zap.Error(err))
	} else {
		log.Info("lnd runninng ",
			zap.String("publicKey", info.IdentityPubkey),
			zap.String("addess", info.IdentityAddress))
	}
	identityPubkey = info.IdentityPubkey
}

type lndHubPrivateServiceServer struct{}

func (lndHubPrivateServiceServer) AddInvoice(ctx context.Context, req *api.AddInvoiceReq) (*empty.Empty, error) {
	userId := ctx.Value(CtxUserIdKey).(string)
	memo := req.Memo
	amount := req.Amount

	// use case
	user := repo.Get(userId)
	resp, err := lnd.AddInvoice(memo, uint(amount))
	user.AttachUserInvoice(resp.PaymentRequest)
	err = repo.Update(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (lndHubPrivateServiceServer) PayInvoice(ctx context.Context, req *api.PayInvoiceReq) (*empty.Empty, error) {
	userId := ctx.Value(CtxUserIdKey).(string)
	invoice := req.Invoice
	amount := req.Amount

	// use case
	user := repo.Get(userId)

	resp, err := lnd.DecodePay(invoice)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	numSatoshis := resp.NumSatoshis
	if numSatoshis == 0 {
		numSatoshis = uint(amount)
	}
	balance := user.Balance()
	if balance >= numSatoshis+uint(math.Floor(float64(numSatoshis)*0.01)) {
		if identityPubkey == resp.Description {
			payee, err := repo.FindByPaymentHash(resp.PaymentHash)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			if user.GetPaymentHashState(resp.PaymentHash) == entity.PAYMENT_HASH_STATE_PAIED {
				return nil, status.Error(codes.Internal, err.Error())
			}

			payee.UpdateBalance(payee.Balance() + numSatoshis)
			err = repo.Update(payee)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			user.UpdateBalance(balance - numSatoshis)
			user.AttachTransaction(
				*entity.NewTx(
					time.Now(),
					"paid_invoice",
					numSatoshis+uint(math.Floor(float64(numSatoshis)*0.01)),
					uint(math.Floor(float64(numSatoshis)*0.03)),
					resp.Description,
				))
			err = repo.Update(user)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			payee.UpdatePaymentHashState(resp.PaymentHash, entity.PAYMENT_HASH_STATE_PAIED)
			return nil, nil
		}
	} else {
		client, err := lnd.GetSendPaymentClient()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		user.UnlockFounds(invoice)
		err = client.Send(lightning.SendRequest{})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return nil, nil
}

func (lndHubPrivateServiceServer) GetBtc(ctx context.Context, req *empty.Empty) (*api.Btc, error) {
	userId := ctx.Value(CtxUserIdKey).(string)

	// use case
	user := repo.Get(userId)
	address := user.Getaddress()

	// ToDo impl: https://github.com/BlueWallet/LndHub/blob/master//class/User.js#L109
	// await this._redis.set('bitcoin_address_for_' + this._userid, address);

	return &api.Btc{
		Address: address,
	}, nil
}

func (lndHubPrivateServiceServer) GetBalance(ctx context.Context, req *empty.Empty) (*api.Balance, error) {
	userId := ctx.Value(CtxUserIdKey).(string)

	// use case
	user := repo.Get(userId)
	return &api.Balance{
		Balance: uint32(user.Balance()),
	}, nil
}

func (lndHubPrivateServiceServer) GetTxs(ctx context.Context, req *empty.Empty) (*api.Transactions, error) {
	panic("implement me")
}

func (lndHubPrivateServiceServer) GetUserInvoices(ctx context.Context, req *empty.Empty) (*api.Invoices, error) {
	userId := ctx.Value(CtxUserIdKey).(string)

	// use case
	user := repo.Get(userId)
	return &api.Invoices{
		Invoice: []string{user.Invoice()},
	}, nil
}

func (lndHubPrivateServiceServer) GetInfo(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	panic("implement me")
}

func (*lndHubPrivateServiceServer) Authorize(ctx context.Context, fullMethodName string) error {
	return nil
}

func GetLndHubPrivateServiceServer() api.LndHubPrivateServiceServer {
	return &lndHubPrivateServiceServer{}
}
