package tester

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/hardware"
	"github.com/ruraomsk/kuda/setup"
)

func C8Tester() {
	for i := 0; i < setup.Set.Hardware.C8count; i++ {
		waitC8(i + 2)
	}
	v := 1
	for {
		for i := 0; i < setup.Set.Hardware.C8count; i++ {
			num := i + 2
			for j := 1; j < 9; j++ {
				time.Sleep(time.Second)
				if hardware.IsWorkC8(num) {
					hardware.C8SetValue(num, j, v)
					if v == 1 {
						time.Sleep(time.Second)
						for k := 0; k < 10; k++ {
							time.Sleep(500 * time.Millisecond)
							hardware.C8SetValue(num, j, 0)
							time.Sleep(500 * time.Millisecond)
							hardware.C8SetValue(num, j, 1)
						}
					}
				} else {
					waitC8(num)

				}
			}
		}
		v++
		if v > 1 {
			v = 0
		}

	}
}

func waitC8(number int) {
	for !hardware.IsWorkC8(number) {
		logger.Info.Printf("wait module c8 %d", number)
		time.Sleep(time.Second)
	}
	logger.Info.Printf("module c8 %d ready", number)
}
