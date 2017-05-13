// Raspberry pi GPIO grpc wrapper

package main

import (
	"fmt"
	pb "github.com/midnightconman/gopio/pb"
	//"github.com/stianeikeland/go-rpio"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"strconv"
)

type server struct{}

func (s *server) GetPinDirection(ctx context.Context, pin *pb.Pin) (*pb.PinDirection, error) {
	//return &pb.PinDirection{number: 1}, nil
	return nil, nil
}

func (s *server) GetPinState(ctx context.Context, pin *pb.Pin) (*pb.PinState, error) {
	return nil, nil
}

func (s *server) GetPinPull(ctx context.Context, pin *pb.Pin) (*pb.PinPull, error) {
	return nil, nil
}

func (s *server) SetPinDirection(ctx context.Context, pin *pb.Pin) (*pb.PinDirection, error) {
	return nil, nil
}

func (s *server) SetPinState(ctx context.Context, pin *pb.Pin) (*pb.PinState, error) {
	return nil, nil
}

func (s *server) SetPinPull(ctx context.Context, pin *pb.Pin) (*pb.PinPull, error) {
	return nil, nil
}

func main() {
	port, err := strconv.ParseInt(os.Getenv("GOPIO_PORT"), 10, 64)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen: %v\n", err)
	}
	var opts []grpc.ServerOption
	if os.Getenv("GOPIO_TLS") != "" {
		creds, err := credentials.NewServerTLSFromFile(os.Getenv("GOPIO_CERT"), os.Getenv("GOPIO_KEY"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate credentials %v\n", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGoPIOServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
