package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MizukiSonoko/lnd-gateway/controller"
	"github.com/MizukiSonoko/lnd-gateway/middleware"
)

func main() {
	log.Printf("start LndHub-go implements")
	errC := make(chan error)
	go func() {
		rootHandler := func(w http.ResponseWriter, r *http.Request) {
			log.Print("root directory is accessed and ignored")
			w.WriteHeader(http.StatusNotFound)
		}
		http.HandleFunc("/", rootHandler)
		for path, h := range controller.GetHandlerFuncs() {
			h := http.HandlerFunc(h)
			// ToDo: umm too bad
			if path != "/login"{
				http.Handle(path, middleware.WithJWT(h))
			}else{
				http.Handle(path, h)
			}
		}
		if err := http.ListenAndServe(":8080", nil); err != nil {
			errC <- err
		}
	}()

	quitC := make(chan os.Signal)
	signal.Notify(quitC, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errC:
		panic(err)
	case <-quitC:
		log.Println("finish!")
	}
}
