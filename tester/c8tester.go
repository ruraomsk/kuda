package tester

import (
	"time"

	"github.com/ruraomsk/kuda/hardware"
	"github.com/ruraomsk/kuda/setup"
)

func cycleL() {
	num := 2
	for {
		time.Sleep(time.Second)
		waitC8(num)
		hardware.C8SetValue(num, 1, 1)
		hardware.C8SetValue(num, 2, 0)
		hardware.C8SetValue(num, 3, 0)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 1, 0)
		hardware.C8SetValue(num, 2, 1)
		hardware.C8SetValue(num, 3, 0)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 1, 0)
		hardware.C8SetValue(num, 2, 0)
		hardware.C8SetValue(num, 3, 1)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 1, 0)
		hardware.C8SetValue(num, 2, 0)
		hardware.C8SetValue(num, 3, 0)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 4, 1)
		hardware.C8SetValue(num, 5, 0)
		hardware.C8SetValue(num, 6, 0)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 4, 0)
		hardware.C8SetValue(num, 5, 1)
		hardware.C8SetValue(num, 6, 0)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 4, 0)
		hardware.C8SetValue(num, 5, 0)
		hardware.C8SetValue(num, 6, 1)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 4, 0)
		hardware.C8SetValue(num, 5, 0)
		hardware.C8SetValue(num, 6, 0)
		time.Sleep(5 * time.Second)
	}

}

func C8Tester() {
	for i := 0; i < setup.Set.Hardware.C8count; i++ {
		num := i + 2
		for j := 1; j < 9; j++ {
			waitC8(num)
			hardware.C8SetValue(num, j, 0)
		}
	}
	time.Sleep(5 * time.Second)
	go cycleL()
	num := 3
	for {
		time.Sleep(time.Second)
		waitC8(num)
		hardware.C8SetValue(num, 1, 1)
		hardware.C8SetValue(num, 2, 0)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 1, 0)
		hardware.C8SetValue(num, 2, 1)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 1, 0)
		hardware.C8SetValue(num, 2, 0)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 3, 1)
		hardware.C8SetValue(num, 4, 0)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 3, 0)
		hardware.C8SetValue(num, 4, 1)
		time.Sleep(5 * time.Second)
		hardware.C8SetValue(num, 3, 0)
		hardware.C8SetValue(num, 4, 0)

	}

}

func waitC8(number int) {
	for !hardware.IsWorkC8(number) {
		// logger.Info.Printf("wait module c8 %d", number)
		time.Sleep(time.Second)
	}
	// logger.Info.Printf("module c8 %d ready", number)
}
