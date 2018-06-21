package main

import (
	"fmt"
	"github.com/midnightconman/gopio/client"
	pb "github.com/midnightconman/gopio/pb"
	"github.com/midnightconman/gopio/schema"
	"io"
	"log"
	"os"
	"time"
)

var (
	Info  *log.Logger
	Error *log.Logger
	state schema.State = schema.Low
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

	conn, err := client.NewClient(server)
	if err != nil {
		Error.Printf("fail to dial: %v", err)
	}
	defer conn.Close()

	pbClient := pb.NewGoPIOClient(conn)

	health, err := client.HealthCheck(pbClient)
	if err != nil {
		Error.Printf("Failed Healthcheck: %v\n", err)
		os.Exit(1)
	}
	Info.Printf("Healthcheck{%v}\n", health)

	// Turn-off all pins before starting
	for _, v := range []int{14, 15, 18, 26, 19, 13, 6, 5} {
		p := pb.Pin{Number: int32(v), Direction: int32(schema.Output), State: int32(schema.High)}
		ps, err := client.PinSet(pbClient, &p)
		if err != nil {
			Error.Printf("Failed PinOutputToggle for pin(%d): %v\n", &p.Number, err)
		}
		Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))
		time.Sleep(time.Duration(1 * int(time.Second)))
	}

	//for k, v := range map[int]int{14: 1, 15: 1, 18: 1, 26: 1, 19: 420, 13: 1, 6: 1, 5: 180} {
	for k, v := range map[int]int{13: 420, 6: 180} {
		p := pb.Pin{Number: int32(k), Direction: int32(schema.Output), State: int32(schema.Low)}
		ps, err := client.PinSet(pbClient, &p)
		if err != nil {
			Error.Printf("Failed PinOutputToggle for pin(%d): %v\n", &p.Number, err)
		}
		Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))

		time.Sleep(time.Duration(v * int(time.Second)))

		p = pb.Pin{Number: int32(k), Direction: int32(schema.Output), State: int32(schema.High)}
		ps, err = client.PinSet(pbClient, &p)
		if err != nil {
			Error.Printf("Failed PinOutputToggle for pin(%d): %v\n", &p.Number, err)
		}
		Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))
	}

}
