

package main

import (
    "flag"
    "fmt"
    "github.com/stianeikeland/go-rpio"
    "net"
    "google.golang.org/grpc"
    pb "github.com/midnightconman/gopio/proto"
)

var (
    tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
    certFile   = flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
    keyFile    = flag.String("key_file", "testdata/server1.key", "The TLS key file")
    port       = flag.Int("port", 8443, "The server port")
)

type server struct{}

const (
    Input Direction = iota
    Output
)

const (
    Low State = iota
    High
)

const (
    PullOff Pull = iota
    PullDown
    PullUp
)

func newServer() *gopioServer {
    s := new(gopioServer)
    return s
}

func main() {
    flag.Parse()
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if err != nil {
        grpclog.Fatalf("failed to listen: %v", err)
    }
    var opts []grpc.ServerOption
    if *tls {
        creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
        if err != nil {
            grpclog.Fatalf("Failed to generate credentials %v", err)
        }
        opts = []grpc.ServerOption{grpc.Creds(creds)}
    }
    grpcServer := grpc.NewServer(opts...)
    //pb.RegisterRouteGuideServer(grpcServer, newServer())
    grpcServer.Serve(lis)
}
