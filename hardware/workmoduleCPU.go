package hardware

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
)

func workModuleCPU() {
	for {
		err := Cpu.setMasterTCP()
		if err != nil {
			logger.Error.Printf("moduleCpu not ready message %s", err.Error())
			time.Sleep(time.Second)
			continue
		}
		break
	}
	Cpu.writer = make(chan writeHR)
	go Cpu.loopTCP()
	for !Cpu.work {
		logger.Info.Printf("moduleCpu wait exchange")
		time.Sleep(time.Second)
	}
	if Cpu.moduleNumber != int(Cpu.masterTCP.hr[Cpu.moduleType]) {
		logger.Error.Printf("moduleCpu not equal %d %d", Cpu.moduleNumber, int(Cpu.masterTCP.hr[Cpu.moduleType]))
		return
	}
	if int(Cpu.masterTCP.hr[Cpu.moduleStatus]) != 1 {
		logger.Error.Printf("moduleCpu status  %d", int(Cpu.masterTCP.hr[Cpu.moduleStatus]))
		return
	}
	for Cpu.work {
		v := 1
		for {
			logger.Info.Printf("moduleCpu work %d", v)
			for i := 1; i < 11; i++ {
				Cpu.SetDO(i, v)
				time.Sleep(time.Second)
			}
			v++
			if v > 1 {
				v = 0
			}

		}
	}
	logger.Info.Printf("moduleCpu ending")
}
