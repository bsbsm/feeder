package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bsbsm/feeder/pkg/db"
	"github.com/bsbsm/feeder/pkg/feeder"
	"github.com/bsbsm/feeder/pkg/server"
)

var readPeriod = flag.Int("rp", 5000, "RSS reading period in ms")

func main() {
	flag.Parse()

	s := db.SQLiteDatabase{}

	server.SetSQLiteDatabase(&s)

	f, err := feeder.NewFeeder(&s)
	if err != nil {
		panic(err)
	}

	go f.Reading(time.Duration(*readPeriod) * time.Millisecond)

	go server.BlockingListen(8080)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-signals
	fmt.Println("\nStop listening")
}
