package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TxCorpi0x/tcp-server/server/srv"
)

func main() {
	s := srv.NewTcpServer(srv.Conf{
		Port:                  "8080",
		ConcurrentConnections: 50,
		ConcurrentHandlers:    3,
	})
	s.Start()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	s.Stop()
}
