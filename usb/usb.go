package usb

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
)

func mainLoop() {
	loop := time.NewTicker(time.Second)
	for {
		select {
		case <-loop.C:

		}
	}
}

func StartUSB() {
	logger.Info.Println("Usb save start")
	go mainLoop()
}
