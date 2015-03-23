package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := Configure(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	closed := false
	listener, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		closed = true
		listener.Close()
	}()

	handler := NewHandler(cfg.Links, !cfg.NoReload, cfg.Allow)
	server := http.Server{Handler: handler}
	if err := server.Serve(listener); err != nil && !closed {
		log.Fatal(err)
	}
}
