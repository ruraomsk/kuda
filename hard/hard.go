package hard

import (
	"time"

	"github.com/ruraomsk/kuda/setup"
)

func mainLoop(stop chan interface{}) {
	loop := time.NewTicker(time.Duration(setup.Set.Hardware.Step) * time.Millisecond)
	for {
		select {
		case <-stop:
			return
		case <-loop.C:

		}
	}
}

func StartHard(stop chan interface{}) {
	go mainLoop(stop)
	WatchDogStart()
}
