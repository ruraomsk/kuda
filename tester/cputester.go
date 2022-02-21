package tester

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/hardware"
)

func CpuTester() {
	waitCpu()
	for hardware.Cpu.IsWork() {
		v := 1
		for {
			logger.Info.Printf("moduleCpu work %d", v)
			for i := 1; i < 11; i++ {
				if hardware.Cpu.IsWork() {
					hardware.Cpu.SetDO(i, v)
					time.Sleep(time.Second)
				} else {
					waitCpu()
				}
			}
			v++
			if v > 1 {
				v = 0
			}

		}
	}

}
func waitCpu() {
	for !hardware.Cpu.IsWork() {
		logger.Info.Println("... wait ready cpu")
		time.Sleep(time.Second)
	}

}
