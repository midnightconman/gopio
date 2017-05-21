package main

import (
	"fmt"
	"github.com/midnightconman/gopio/client"
	pb "github.com/midnightconman/gopio/pb"
	"github.com/midnightconman/gopio/schema"
	"gopkg.in/robfig/cron.v2"
	"io"
	"log"
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

func SignalHandler() bool {
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

	return true
}

func main() {
	LogInit(os.Stdout, os.Stderr)
	SignalHandler()
	server := fmt.Sprintf("%s:%s", os.Getenv("GOPIO_HOST"), os.Getenv("GOPIO_PORT"))

	//var opts []grpc.DialOption
	//conn, err := client.NewClient(server, opts...)
	conn, err := client.NewClient(server)
	if err != nil {
		Error.Printf("fail to dial: %v", err)
	}
	defer conn.Close()

	pbClient := pb.NewGoPIOClient(conn)

	c := cron.New()

	p := pb.Pin{Number: 14, Direction: int32(schema.Output), State: int32(schema.Low)}
	c.AddFunc(os.Getenv("GOPIO_CRON_ON"), func() {
		ps, err := client.PinSet(pbClient, &p)
		if err != nil {
			Error.Printf("Failed PinOn for pin(%d): %v\n", &p.Number, err)
		}
		Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))
	})

	p = pb.Pin{Number: 14, Direction: int32(schema.Output), State: int32(schema.High)}
	c.AddFunc(os.Getenv("GOPIO_CRON_OFF"), func() {
		ps, err := client.PinSet(pbClient, &p)
		if err != nil {
			Error.Printf("Failed PinOn for pin(%d): %v\n", &p.Number, err)
		}
		Info.Printf("Pin:(%d) Direction:(%s) State:(%s)\n", p.Number, schema.Direction(uint8(ps.Direction)), schema.State(uint8(ps.State)))
	})

	c.Start()

	//for {
	//}

}
