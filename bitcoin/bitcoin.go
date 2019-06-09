package bitcoin

import (
	"github.com/MizukiSonoko/LndHub-go/logger"
	"github.com/btcsuite/btcd/rpcclient"
	"go.uber.org/zap"
)

var (
	log = logger.NewLogger()
	bc  = new(bClient)
)

type BitcoinClient interface {
	ImportAddress(address string) error
}

func init() {
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "sonokko",
		Pass:         "1qazxsw23edcvfr4 ",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal("initialize bitcoin client failed", zap.Error(err))
	}

	if err := client.Ping(); err != nil {
		log.Fatal("bitcoin client Ping failed", zap.Error(err))
	}
	bc.client = client
}

type bClient struct {
	client *rpcclient.Client
}

func (b *bClient) ImportAddress(address string) error {
	return b.client.ImportAddress(address)
}

func GetBitcoinClient() BitcoinClient {
	return bc
}
