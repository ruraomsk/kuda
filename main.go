package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/hardware"
	"github.com/ruraomsk/kuda/netware"
	"github.com/ruraomsk/kuda/setup"
	"github.com/ruraomsk/kuda/status"
	"github.com/ruraomsk/kuda/tester"
	"github.com/ruraomsk/kuda/transport"
	"github.com/ruraomsk/kuda/usb"
)

var (
	//go:embed config
	config embed.FS
)

func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFS(config, "config/config.toml", &setup.Set); err != nil {
		fmt.Println("Dissmis config.toml")
		os.Exit(-1)
		return
	}
	os.MkdirAll(setup.Set.LogPath, 0777)
	os.MkdirAll(setup.Set.SetupBrams.DbPath, 0777)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := logger.Init(setup.Set.LogPath); err != nil {
		log.Panic("Error logger system", err.Error())
		return
	}
	dbstop := make(chan interface{})

	brams.StartBrams(dbstop)

	if err := status.StartStatus(); err != nil {
		logger.Error.Printf("Subsystem status %s", err.Error())
		return
	}
	hardware.StartHard()
	netware.StartNetware()
	usb.StartUSB()
	go transport.StartServerExchange("192.168.115.159:2018")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Println("kuda start")
	logger.Info.Println("kuda start")
	watch := time.NewTicker(time.Duration(setup.Set.WatchDog.Step) * time.Millisecond)
	go tester.CpuTester()
	// go tester.C8Tester()
	//tester.BinTest()
	buffer, err := config.ReadFile("config/test.json")
	if err != nil {
		logger.Error.Printf("test.json %s", err.Error())
		fmt.Println(err.Error())
		return
	}
	go tester.RpuTest(buffer)

loop:
	for {
		select {
		case <-c:
			fmt.Println("Wait make abort...")
			hardware.ExitV220()
			if transport.IsConnected {
				transport.ExitDevice <- 1
			}
			time.Sleep(3 * time.Second)
			dbstop <- 1
			// hardstop <- 1
			time.Sleep(3 * time.Second)
			break loop
		case <-watch.C:
			hardware.WatchDogTick()

		}
	}
	fmt.Println("kuda stop")
	logger.Info.Println("kuda stop")

}
