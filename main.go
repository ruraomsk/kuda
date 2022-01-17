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
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/hard"
	"github.com/ruraomsk/kuda/netware"
	"github.com/ruraomsk/kuda/setup"
	"github.com/ruraomsk/kuda/status"
	"github.com/ruraomsk/kuda/usb"
)

var (
	//go:embed config
	config embed.FS
)

func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFS(config, "config/config.toml", &setup.Set); err != nil {
		fmt.Println("Отсутствует config.toml")
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
	hardstop := make(chan interface{})

	brams.StartBrams(dbstop)

	if err := status.StartStatus(); err != nil {
		logger.Error.Printf("Подсистема status %s", err.Error())
		return
	}
	hard.StartHard(hardstop)
	netware.StartNetware()
	usb.StartUSB()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	fmt.Println("kuda start")
	logger.Info.Println("kuda start")
	watch := time.NewTicker(time.Duration(setup.Set.WatchDog.Step) * time.Millisecond)
	for {
		select {
		case <-c:
			fmt.Println("Wait make abort...")
			dbstop <- 1
			hardstop <- 1
			time.Sleep(3 * time.Second)
			fmt.Println("kuda stop")
			logger.Info.Println("kuda stop")
			os.Exit(0)

		case <-watch.C:
			hard.WatchDogTick()
		}
	}

}
