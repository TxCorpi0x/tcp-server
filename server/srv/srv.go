package srv

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type connectionPool chan net.Conn

type semaphore chan int

type Conf struct {
	Port string

	ConcurrentConnections int
	ConcurrentHandlers    int
}

type tcpServer struct {
	conf Conf

	listener net.Listener

	quit chan bool

	connectionPool     connectionPool
	concurrentHandlers semaphore

	// mutex of num concurrency protection
	mu  sync.Mutex
	num uint64
}

func NewTcpServer(conf Conf) *tcpServer {
	if conf.ConcurrentConnections == 0 {
		fmt.Println("Error: Minimum concurrent connections should be more than zero")
		return nil
	}
	if conf.ConcurrentHandlers == 0 {
		fmt.Println("Error: Minimum concurrent handlers should be more than zero")
		return nil
	}

	var err error
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", conf.Port))
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	s := tcpServer{
		conf:               conf,
		listener:           listener,
		quit:               make(chan bool, 1),
		concurrentHandlers: make(chan int, conf.ConcurrentHandlers),
		connectionPool:     make(chan net.Conn, conf.ConcurrentConnections),
		num:                0,
	}
	return &s
}

// accepts new connection and put them into the connection pool channel
func (s *tcpServer) accept() {
	for {
		select {
		case <-s.quit:
			fmt.Println("Accept is returning warnings to the rest of client requests...")
			for {
				conn, err := s.listener.Accept()
				if err != nil {
					continue
				}
				conn.Write([]byte("Warning: Server will not add this request, the listener is quitting"))
			}
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}

			// add new connection to the pool
			s.connectionPool <- conn
		}
	}
}

// handles all connections concurrently
func (s *tcpServer) handle() {
	for {
		conn := <-s.connectionPool
		// acquire concurrent handlers semaphore capacity
		s.concurrentHandlers <- 1
		go func() {
			s.handleConnection(conn)

			// // TODO: uncomment to simulate time consuming process
			// sec := time.Duration(rand.Intn(5-2) + 2)
			// time.Sleep(sec * time.Second)

			// release concurrent handlers semaphore capacity
			<-s.concurrentHandlers
		}()
	}
}

func (s *tcpServer) Start() {
	go s.accept()
	go s.handle()
}

func (s *tcpServer) handleConnection(conn net.Conn) {
	defer func() {
		s.mu.Unlock()
		conn.Close()
	}()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read error", err.Error())
		return
	}

	s.mu.Lock()
	s.num++

	// Write a response back to the client.
	conn.Write([]byte(strconv.FormatUint(s.num, 10)))
}

func (s *tcpServer) waitForClose() {
	fmt.Println("Waiting for open precesses...")
	for {
		if len(s.connectionPool) == 0 && len(s.concurrentHandlers) == 0 {
			fmt.Printf("Num: %d, Open Connections: %d, Open Processes: %d\n", s.num, len(s.connectionPool), len(s.concurrentHandlers))
			fmt.Println("All Processes Done!")

			return
		}
	}
}

func (s *tcpServer) Stop() {
	close(s.quit)
	fmt.Println("Shutting down...")

	allProcessed := make(chan struct{})
	go func() {
		s.waitForClose()
		close(allProcessed)
	}()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Printf("Num: %d, Open Connections: %d, Open Processes: %d\n", s.num, len(s.connectionPool), len(s.concurrentHandlers))
		}
	}()

	<-allProcessed
	s.listener.Close()
}
