package cli

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	Conf Conf
}

type Conf struct {
	Port string
}

func NewClient(conf Conf) *Client {
	return &Client{
		Conf: conf,
	}
}

var num int = 0

func (c *Client) Send() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%s", c.Conf.Port))
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	num++
	numToSend := num
	_, err = conn.Write([]byte("Add"))
	if err != nil {
		fmt.Printf("Write to server failed for %d: %s\n", numToSend, err.Error())
		os.Exit(1)
	}

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		println("Read from server reply failed:", err.Error())
		os.Exit(1)
	}

	if strings.HasPrefix(string(reply), "Warning") {
		fmt.Printf("reply from server for %d= %s\n", numToSend, string(reply))
	} else {
		fmt.Printf("reply from server= %s\n", string(reply))
	}

	conn.Close()
}
