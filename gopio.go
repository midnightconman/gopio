// Raspberry pi GPIO grpc wrapper

package main

import (
	"fmt"
	pb "github.com/midnightconman/gopio/pb"
	"github.com/midnightconman/gopio/schema"
	"github.com/stianeikeland/go-rpio"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func LogInit(
	infoHandle io.Writer,
	errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

type server struct{}

func (s *server) GetPinDirection(ctx context.Context, pin *pb.Pin) (*pb.PinDirection, error) {
	/*p := rpio.Pin(pin.Number)
	if err := rpio.Open(); err != nil {
		return &pb.PinDirection{Direction: Input}, err
	}

	defer rpio.Close()

	Info.Printf("GetPinDirection context(%s) pin(%v)\n", ctx, pin)

	// This isn't supported by rpio lib yet
	return &pb.PinDirection{Direction: Input}, nil
	*/
	return nil, nil
}

func (s *server) GetPinState(ctx context.Context, pin *pb.Pin) (*pb.PinState, error) {
	p := rpio.Pin(pin.Number)
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	defer rpio.Close()

	Info.Printf("TogglePinState context(%s) pin{Number: %d, Direction: %s, State: %s, Pull: %s}\n", ctx, pin.Number, schema.Direction(int32(pin.Direction)), schema.State(int32(pin.State)), schema.Pull(int32(pin.Pull)))

	return &pb.PinState{State: int32(p.Read())}, nil
}

func (s *server) GetPinPull(ctx context.Context, pin *pb.Pin) (*pb.PinPull, error) {
	/*p := rpio.Pin(pin.Number)
	if err := rpio.Open(); err != nil {
		return &pb.PinPull{Pull: p.GetPull()}, err
	}

	defer rpio.Close()

	Info.Printf("TogglePinState context(%s) pin{Number: %d, Direction: %s, State: %s, Pull: %s}\n", ctx, pin.Number, schema.Direction(int32(pin.Direction)), schema.State(int32(pin.State)), schema.Pull(int32(pin.Pull)))

	// This isn't supported by rpio lib yet
	return &pb.PinPull{Pull: int32(p.GetPinPull())}, nil
	*/
	return nil, nil
}

func (s *server) SetPinDirection(ctx context.Context, pin *pb.Pin) (*pb.PinDirection, error) {
	p := rpio.Pin(pin.Number)
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	defer rpio.Close()
	p.Mode(rpio.Direction(uint8(pin.Direction)))

	Info.Printf("TogglePinState context(%s) pin{Number: %d, Direction: %s, State: %s, Pull: %s}\n", ctx, pin.Number, schema.Direction(int32(pin.Direction)), schema.State(int32(pin.State)), schema.Pull(int32(pin.Pull)))

	return &pb.PinDirection{Direction: pin.Direction}, nil
}

func (s *server) SetPinState(ctx context.Context, pin *pb.Pin) (*pb.PinState, error) {
	p := rpio.Pin(pin.Number)
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	defer rpio.Close()
	p.Write(rpio.State(uint8(pin.State)))

	Info.Printf("TogglePinState context(%s) pin{Number: %d, Direction: %s, State: %s, Pull: %s}\n", ctx, pin.Number, schema.Direction(int32(pin.Direction)), schema.State(int32(pin.State)), schema.Pull(int32(pin.Pull)))

	return &pb.PinState{State: int32(p.Read())}, nil
}

func (s *server) SetPinPull(ctx context.Context, pin *pb.Pin) (*pb.PinPull, error) {
	p := rpio.Pin(pin.Number)
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	defer rpio.Close()
	p.Write(rpio.State(uint8(pin.State)))

	Info.Printf("TogglePinState context(%s) pin{Number: %d, Direction: %s, State: %s, Pull: %s}\n", ctx, pin.Number, schema.Direction(int32(pin.Direction)), schema.State(int32(pin.State)), schema.Pull(int32(pin.Pull)))

	return &pb.PinPull{Pull: pin.Pull}, nil
}

func (s *server) TogglePinState(ctx context.Context, pin *pb.Pin) (*pb.PinState, error) {
	p := rpio.Pin(pin.Number)
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	defer rpio.Close()
	p.Toggle()

	Info.Printf("TogglePinState context(%s) pin{Number: %d, Direction: %s, State: %s, Pull: %s}\n", ctx, pin.Number, schema.Direction(int32(pin.Direction)), schema.State(int32(pin.State)), schema.Pull(int32(pin.Pull)))

	return &pb.PinState{State: int32(p.Read())}, nil
}

func main() {
	LogInit(os.Stdout, os.Stderr)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sig := <-sigs
			Info.Printf("Signal Received: %v\n", sig)
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				os.Exit(0)
			}
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Getenv("GOPIO_HOST"), os.Getenv("GOPIO_PORT")))
	if err != nil {
		Error.Printf("failed to listen: %v\n", err)
	}
	var opts []grpc.ServerOption
	if os.Getenv("GOPIO_TLS") != "" {
		creds, err := credentials.NewServerTLSFromFile(os.Getenv("GOPIO_CERT"), os.Getenv("GOPIO_KEY"))
		if err != nil {
			Error.Printf("failed to generate credentials %v\n", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGoPIOServer(grpcServer, &server{})
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)
}
