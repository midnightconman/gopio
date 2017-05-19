package main

import (
	"fmt"
	pb "github.com/midnightconman/gopio/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"os"
	"time"
)

type Direction uint8
type State uint8
type Pull uint8

const (
	Input Direction = iota
	Output
)
const _Direction_name = "InputOutput"

var _Direction_index = [...]uint8{0, 5, 11}

func (i Direction) String() string {
	if i >= Direction(len(_Direction_index)-1) {
		return fmt.Sprintf("Direction(%d)", i)
	}
	return _Direction_name[_Direction_index[i]:_Direction_index[i+1]]
}

const (
	Low State = iota
	High
)
const _State_name = "LowHigh"

var _State_index = [...]uint8{0, 3, 7}

func (i State) String() string {
	if i >= State(len(_State_index)-1) {
		return fmt.Sprintf("State(%d)", i)
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}

const (
	PullOff Pull = iota
	PullDown
	PullUp
)
const _Pull_name = "PullOffPullDownPullUp"

var _Pull_index = [...]uint8{0, 7, 15, 21}

func (i Pull) String() string {
	if i >= Pull(len(_Pull_index)-1) {
		return fmt.Sprintf("Pull(%d)", i)
	}
	return _Pull_name[_Pull_index[i]:_Pull_index[i+1]]
}

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

	p := pb.Pin{Number: 14, Direction: int32(Output), State: int32(Low)}
	ps, err := PinSet(client, &p)
	if err != nil {
		Error.Printf("Failed PinOn for pin(%d): %v\n", &p.Number, err)
	}
	Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, Direction(uint8(ps.Direction)), State(uint8(ps.State)))

	time.Sleep(10 * time.Second)

	p = pb.Pin{Number: 14, Direction: int32(Output), State: int32(High)}
	ps, err = PinSet(client, &p)
	if err != nil {
		Error.Printf("Failed PinOn for pin(%d): %v\n", &p.Number, err)
	}
	Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, Direction(uint8(ps.Direction)), State(uint8(ps.State)))

}
