package tester

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/hardware"
	"github.com/ruraomsk/kuda/hardware/bin"
	"github.com/ruraomsk/kuda/setup"
)

var cmk bin.CMK

func RpuTest(b []byte) {

	err := json.Unmarshal(b, &cmk)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	for !hardware.AllReady() {
		fmt.Println("wait...")
		time.Sleep(time.Second)
	}
	//Обнуляем все тиристоры
	zeroOn()

osStop:
	osWork()
allBlink:
	allYellowBlink()
	zeroOn()
	//Выдаем кругом красный
	for _, v := range cmk.TirToNaps {
		for _, d := range v.Reds {
			hardware.C8SetOut(d, 1)
		}
	}
	time.Sleep(time.Duration(time.Duration(setup.Set.Hardware.LongKK) * time.Second))
	zeroOn()
	//Стартуем первую фазу
	control := time.NewTicker(100 * time.Millisecond)
	ch := 10
	changePhase := time.NewTimer(time.Duration(ch) * time.Second)
	for {
		select {
		case <-control.C:
			if hardware.Cpu.GetDI(setup.Set.Hardware.PinOS) {
				control.Stop()
				goto osStop
			}
			if hardware.Cpu.GetDI(setup.Set.Hardware.PinYB) {
				control.Stop()
				goto allBlink
			}
		case <-changePhase.C:
			fmt.Printf("change phase %d\n", ch)
			ch++
			changePhase = time.NewTimer(time.Duration(ch) * time.Second)
		}

	}

}
func osWork() bool {
	result := false
	for hardware.Cpu.GetDI(setup.Set.Hardware.PinOS) {
		result = true
		zeroOn()
		time.Sleep(100 * time.Millisecond)
	}
	return result
}
func allYellowBlink() {
	for hardware.Cpu.GetDI(setup.Set.Hardware.PinYB) {
		if osWork() {
			continue
		}
		zeroOn()
		time.Sleep(500 * time.Millisecond)
		if osWork() {
			continue
		}
		for _, v := range cmk.TirToNaps {
			if v.Yellow > 0 {
				hardware.C8SetOut(v.Yellow, 1)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	zeroOn()
}
func zeroOn() {
	for i := 0; i < setup.Set.Hardware.C8count; i++ {
		num := i + 2
		for j := 1; j < 9; j++ {
			waitC8(num)
			hardware.C8SetValue(num, j, 0)
		}
	}
}
