package main

import (
	"fmt"
	"github.com/midnightconman/gopio/client"
	pb "github.com/midnightconman/gopio/pb"
	"github.com/midnightconman/gopio/schema"
	"io"
	"log"
	"os"
	"strconv"
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

	for _, v := range []int{14, 15, 18, 26, 19, 13, 6, 5} {
		for i := 0; i < 2; i++ {
			if i%2 == 0 {
				state = schema.Low
			} else {
				state = schema.High
			}
			p := pb.Pin{Number: int32(v), Direction: int32(schema.Output), State: int32(state)}
			ps, err := client.PinSet(pbClient, &p)
			if err != nil {
				Error.Printf("Failed PinOutputToggle for pin(%d): %v\n", &p.Number, err)
			}
			Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))

			if os.Getenv("GOPIO_DELAY") != "" {
				i, _ := strconv.Atoi(os.Getenv("GOPIO_DELAY"))
				time.Sleep(time.Duration(i * int(time.Second)))
			} else {
				time.Sleep(10 * time.Second)
			}
		}
	}

}
