package main

import (
	"fmt"
    pb "github.com/midnightconman/gopio/pb"
	"github.com/midnightconman/gopio/client"
	"github.com/midnightconman/gopio/schema"
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

func main() {
	LogInit(os.Stdout, os.Stderr)
	server := fmt.Sprintf("%s:%s", os.Getenv("GOPIO_HOST"), os.Getenv("GOPIO_PORT"))

	//var opts []grpc.DialOption

	//conn, err := client.NewClient(server, opts...)
	conn, err := client.NewClient(server)
	if err != nil {
		Error.Printf("fail to dial: %v", err)
	}
	defer conn.Close()

	pbClient := pb.NewGoPIOClient(conn)

    //for i := 2; i <= 27; i++ {

	p := pb.Pin{Number: 14, Direction: int32(schema.Output), State: int32(schema.Low)}
	ps, err := client.PinSet(pbClient, &p)
	if err != nil {
		Error.Printf("Failed PinOn for pin(%d): %v\n", &p.Number, err)
	}
	Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))

	time.Sleep(1 * time.Second)

	p = pb.Pin{Number: 14, Direction: int32(schema.Output), State: int32(schema.High)}
	ps, err = client.PinSet(pbClient, &p)
	if err != nil {
		Error.Printf("Failed PinOn for pin(%d): %v\n", &p.Number, err)
	}
	Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))
    //}

}
