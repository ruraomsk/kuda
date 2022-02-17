package hardware

import (
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/kuda/brams"
	"github.com/ruraomsk/kuda/status"
)

// func (m *MasterRtu) readAllHR() error {
// 	if len(m.hr) == 0 {
// 		return nil
// 	}
// 	size := uint16(100)
// 	ref := uint16(0)
// 	for count := len(m.hr); count > 0; count -= int(size) {
// 		len := count
// 		if count > int(size) {
// 			len = int(size)
// 		}
// 		buff, err := m.client.ReadHoldingRegisters(ref, uint16(len))
// 		if err != nil {
// 			logger.Error.Printf("read hr %d %d %s", ref, len, err.Error())
// 			return err
// 		}
// 		pos := ref
// 		left := 0
// 		for i := 0; i < len; i++ {
// 			m.hr[pos] = (uint16(buff[left]) << 8) | uint16(buff[left+1])
// 			pos++
// 			left += 2
// 		}
// 		ref += size
// 	}
// 	return nil
// }

// func mainLoopRTU(stop chan interface{}) {
// 	m := MasterRtu{hr: make([]uint16, 549)}
// 	m.master = *modbus.NewRTUClientHandler("/dev/ttyS1")
// 	m.master.SlaveId = 255
// 	m.master.Logger = logger.Error
// 	m.master.Timeout = time.Second
// 	m.master.IdleTimeout = time.Minute
// 	m.master.BaudRate = 115200
// 	m.master.DataBits = 8
// 	m.master.StopBits = 1
// 	m.master.Parity = "N"
// 	m.master.RS485 = serial.RS485Config{Enabled: true}

// 	if err := m.master.Connect(); err != nil {
// 		logger.Error.Printf("error modbus %s", err.Error())
// 		return
// 	}
// 	m.client = modbus.NewClient(&m.master)

// 	loop := time.NewTicker(time.Duration(setup.Set.Hardware.Step) * time.Millisecond)
// 	for {
// 		select {
// 		case <-stop:
// 			return
// 		case <-loop.C:
// 			if err := m.readAllHR(); err == nil {
// 				logger.Info.Printf("%d %v", m.master.SlaveId, m.hr)
// 			}
// 		}
// 	}
// }

func StartHard() {
	go workModuleCPU()
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
