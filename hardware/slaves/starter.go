package slaves

import (
	"time"

	"github.com/ruraomsk/kuda/setup"
)

func StartDeviceC8() {
	port := setup.Set.Hardware.SPort
	devs := make(map[int]*DeviceC8)
	for i := 0; i < setup.Set.Hardware.C8count; i++ {
		d := new(DeviceC8)
		d.Port = port
		d.Number = i + 2
		d.Size = 6
		devs[d.Number] = d
		go d.WorkDevice()
		port++
	}
	for {
		time.Sleep(time.Second)

	}
}
