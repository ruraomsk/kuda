package bin

import (
	"fmt"
	"time"

	"github.com/ruraomsk/kuda/modbus"
	"github.com/ruraomsk/kuda/setup"
)

var counts map[int]*Counter
var ids []int
var AU = 999
var modbusWork = false
var chanIDS chan []int

func counterStart() {
	counts = make(map[int]*Counter)
	chanIDS = make(chan []int, 100)
	maxId := 0
	for _, v := range cmk.Counters {
		counts[v.Number] = &Counter{Number: v.Number, ID: v.ID, Type: v.Type, Value: v.Default, Default: v.Default}
		if v.ID > maxId {
			maxId = v.ID
		}
	}
	go modbusExch()
	ids = make([]int, maxId)
}
func resetCounters() {
	for _, v := range counts {
		v.Value = v.Default
		v.NeedSend = false
	}
}
func setCounter(number, value int) {
	c, is := counts[number]
	if !is {
		return
	}
	c.Value = value
	c.NeedSend = true
}
func oneStepCounters() {
	needSendMore := false
	for _, v := range counts {
		if !v.NeedSend {
			continue
		}
		needSendMore = true
		ids[v.ID] = v.Value
		if v.Type == 0 {
			if v.Value != AU {
				v.Value--
				if v.Value < 0 {
					v.Value = 0
					v.NeedSend = false
				}
			} else {
				v.NeedSend = false
			}

		} else {
			v.NeedSend = false
		}
	}
	if needSendMore {
		chanIDS <- ids
	}
}

func modbusExch() {
	con := fmt.Sprintf("%s:%d", setup.Set.Counter.Connect, setup.Set.Counter.SPort)
	master := modbus.NewTCPClientHandler(con)
	master.SlaveId = byte(1)
	master.Timeout = time.Second
	master.IdleTimeout = time.Minute
	for {
		modbusWork = false
		if err := master.Connect(); err != nil {

			fmt.Printf("error modbus %s %s", con, err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		client := modbus.NewClient(master)
		modbusWork = true
	work:
		for {
			select {
			case ids := <-chanIDS:
				buf := make([]byte, len(ids)*2)
				for i := 0; i < len(ids); i++ {
					buf[i*2] = byte((ids[i] >> 8))
					buf[i*2+1] = byte((ids[i] & 0xff))

				}
				if _, err := client.WriteMultipleRegisters(0, uint16(len(ids)), buf); err != nil {
					fmt.Printf("error modbus %s %s", con, err.Error())
					modbusWork = false
					break work
				}
			}

		}
		master.Close()
	}

}
