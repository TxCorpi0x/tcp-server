# Asynchronous TCP Server

The present repo, implements a TCP server which is able to process the incoming TCP requests and reply them in asynchronously.

## Features

1. Separate concurrent configurable connection pool.
2. Separate concurrent request processor.
3. Graceful shutdown with connection pool and running process control.

## Usage

### Configurations

Connection handlers, asynchronously process the connections of the connection pool. so the processes will randomly pick the connections. this means the results may return with different order than requests.  

- Port: TCP server listen port
- ConcurrentConnections: Connection pool concurrent connection count.
- ConcurrentHandlers: Number of Concurrent processes that process the connection pool requests asynchronously.

> Note: To increase the performance of the server we can increase the concurrent connections and concurrent handlers (e.g. 100, 200) which means 100 connections will be processed by 200 processors.

```go
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
```

or

```bash
cd srv
go run main.go
```

## Client

A sample concurrent code requests will be send to the server's endpoint.

```go
cd client
go run main.go
```

## Simulation

There is a `TODO` line in the server implementation, uncomment to simulate a random time consuming process.

```go
// TODO: uncomment to simulate time consuming process
```
