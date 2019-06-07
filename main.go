package main

import (
	"context"
	"github.com/MizukiSonoko/lnd-gateway/controller"
	"github.com/MizukiSonoko/lnd-gateway/logger"
	"github.com/MizukiSonoko/lnd-gateway/protobuf"
	"os"
	"os/signal"
	"syscall"
)

var (
	log = logger.NewLogger()
)

func main() {
	log.Info("start LndHub-go server")
	errC := make(chan error)

	go func() {
		ctx := context.Background()
		s := controller.NewGRPCServer(ctx)

		api.RegisterLndHubServiceServer(s.Server, controller.GetLndHubServiceServer())

		api.RegisterLndHubPrivateServiceServer(s.Server, controller.GetLndHubPrivateServiceServer())

		if err := s.Start(); err != nil {
			errC <- err
		}
	}()

	quitC := make(chan os.Signal, 1)
	signal.Notify(quitC, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errC:
		panic(err)
	case <-quitC:
		log.Info("system shutdown")
	}
}
