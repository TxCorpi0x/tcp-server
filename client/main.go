package main

import (
	"sync"

	"github.com/TxCorpi0x/tcp-server/client/cli"
)

func main() {
	c := cli.NewClient(cli.Conf{
		Port: "8080",
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		go c.Send()
		wg.Add(1)
	}
	wg.Wait()
}
