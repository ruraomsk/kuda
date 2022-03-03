package hardware

import (
	"fmt"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/setup"
	"github.com/ruraomsk/kuda/status"
)

var MapC8 map[int]*ModuleC8

func AllReady() bool {
	for _, v := range MapC8 {
		if !v.work {
			return false
		}
	}
	return Cpu.IsWork()
}
func StartHard() {
	MapC8 = make(map[int]*ModuleC8)
	go workModuleCPU()
	port := setup.Set.Hardware.SPort
	for i := 0; i < setup.Set.Hardware.C8count; i++ {
		m := new(ModuleC8)
		m.с8 = C8s
		m.connect = fmt.Sprintf("%s:%d", setup.Set.Hardware.Connect, port)
		m.moduleNumber = i + 2
		m.moduleSlaveID = i + 1
		m.moduleType = 0
		m.moduleStatus = 1
		m.moduleSetup = 2
		m.size = 7
		if err := m.setMasterTCP(); err != nil {
			logger.Error.Printf("%s %s", m.connect, err.Error())
		} else {
			MapC8[m.moduleNumber] = m
			go m.loopTCP()
			logger.Info.Printf("start modbus %s", m.connect)
		}
		port++
	}

	// go mainLoopRTU(stop)
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
