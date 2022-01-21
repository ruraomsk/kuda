package hard

import (
	"time"

	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/setup"
	"github.com/ruraomsk/kuda/status"
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
	status.HardMessage("Запущено оборудование")
	WatchDogStart()
}

func ExitV220() {
	var ed pudge.ErrorDevice
	db, err := brams.Open("base")
	if err != nil {
		return
	}
	db.ReadJSON(&ed)
	ed.V220DK1 = true
	db.WriteJSON(ed)
	db.Close()
}
