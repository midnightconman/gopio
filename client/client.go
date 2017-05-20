package main

import (
	"fmt"
	pb "github.com/midnightconman/gopio/pb"
	"github.com/midnightconman/gopio/schema"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"os"
	"time"
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

func PinSet(client pb.GoPIOClient, pin *pb.Pin) (*pb.Pin, error) {

	d, err := client.SetPinDirection(context.Background(), pin)
	if err != nil {
		return &pb.Pin{Number: 14}, fmt.Errorf("Failed SetPinDirection for pin(%d): %v\n", pin.Number, err)
	}

	s, err := client.SetPinState(context.Background(), pin)
	if err != nil {
		return &pb.Pin{Number: 14}, fmt.Errorf("Failed SetPinState for pin(%d): %v\n", pin.Number, err)
	}

	return &pb.Pin{Number: pin.Number, Direction: int32(d.Direction), State: int32(s.State)}, nil
}

func main() {
	LogInit(os.Stdout, os.Stderr)
	server := fmt.Sprintf("%s:%s", os.Getenv("GOPIO_HOST"), os.Getenv("GOPIO_PORT"))

	var opts []grpc.DialOption
	if os.Getenv("GOPIO_TLS") != "" {
		creds, err := credentials.NewServerTLSFromFile(os.Getenv("GOPIO_CERT"), os.Getenv("GOPIO_KEY"))
		if err != nil {
			Error.Printf("failed to generate credentials %v\n", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(server, opts...)
	if err != nil {
		Error.Printf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewGoPIOClient(conn)

	//for i := 2; i <= 27; i++ {

	p := pb.Pin{Number: 14, Direction: int32(schema.Output), State: int32(schema.Low)}
	ps, err := PinSet(client, &p)
	if err != nil {
		Error.Printf("Failed PinOn for pin(%d): %v\n", &p.Number, err)
	}
	Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))

	time.Sleep(1 * time.Second)

	//p = pb.Pin{Number: int32(i), Direction: int32(schema.Output), State: int32(schema.High)}
	p = pb.Pin{Number: 14, Direction: int32(schema.Output), State: int32(schema.High)}
	ps, err = PinSet(client, &p)
	if err != nil {
		Error.Printf("Failed PinOn for pin(%d): %v\n", &p.Number, err)
	}
	Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))
	//}

}
