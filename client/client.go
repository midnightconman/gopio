package client

import (
	"fmt"
	pb "github.com/midnightconman/gopio/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewClient(server string) (*grpc.ClientConn, error) {

	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("fail to dial: %v", err)
	}

	return conn, nil
}

func NewTLSClient(server string, cert string, key string) (*grpc.ClientConn, error) {

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate credentials %v\n", err)
	}

	conn, err := grpc.Dial(server, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("fail to dial: %v", err)
	}

	return conn, nil
}

func HealthCheck(ctx context.Context, client pb.GoPIOClient) (*pb.Health, error) {
	h, err := client.HealthCheck(ctx, &pb.Health{Alive: true})
	if err != nil {
		return &pb.Health{Alive: false}, fmt.Errorf("Failed Healthcheck: %v\n", err)
	}

	return h, nil
}

func PinSet(ctx context.Context, client pb.GoPIOClient, pin *pb.Pin) (*pb.Pin, error) {

	d, err := client.SetPinDirection(ctx, pin)
	if err != nil {
		return &pb.Pin{Number: 14}, fmt.Errorf("Failed SetPinDirection for pin(%d): %v\n", pin.Number, err)
	}

	s, err := client.SetPinState(ctx, pin)
	if err != nil {
		return &pb.Pin{Number: 14}, fmt.Errorf("Failed SetPinState for pin(%d): %v\n", pin.Number, err)
	}

	return &pb.Pin{Number: pin.Number, Direction: int32(d.Direction), State: int32(s.State)}, nil
}
