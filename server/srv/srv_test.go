package srv_test

import (
	"testing"
	"time"

	"github.com/TxCorpi0x/tcp-server/server/srv"
)

func TestServer(t *testing.T) {
	s := srv.NewTcpServer(srv.Conf{
		Port:                  "8080",
		ConcurrentConnections: 50,
		ConcurrentHandlers:    3,
	})
	s.Start()
	time.Sleep(2 * time.Second)
	s.Stop()
}
