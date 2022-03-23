package bin

import (
	"fmt"
	"time"

	"github.com/ruraomsk/kuda/modbus"
	"github.com/ruraomsk/kuda/setup"
)

var counts map[int][]Counter
var AU = 999
var modbusWork = false
var chanIDS chan Counter
var decoder map[int]Counter

func counterStart() {
	decoder = make(map[int]Counter)
	counts = make(map[int][]Counter)
	chanIDS = make(chan Counter, 100)
	maxId := 0
	for _, v := range cmk.Counters {
		decoder[v.Number] = v
		if v.ID > maxId {
			maxId = v.ID
		}
	}
	go modbusExch()
}
func resetCounters() {
	counts = make(map[int][]Counter)
	for _, v := range decoder {
		if modbusWork {
			chanIDS <- Counter{ID: v.ID, Value: 0}
		}
	}
}

func setCounter(time, number, value int) {
	if time >= 86400 {
		time -= 86400
	}
	c, is := counts[time]
	if !is {
		c = make([]Counter, 0)
	}
	w := decoder[number]
	found := false
	for _, v := range c {
		if v.ID == w.ID && w.Value == v.Value {
			//Дубликат
			found = true
		}
	}
	if !found {
		c = append(c, Counter{ID: w.ID, Value: value})
		counts[time] = c
	}
}
func oneStepCounters() {
	sc := TimeNowOfSecond()
	c, is := counts[sc]
	if !is {
		return
	}
	fmt.Printf("time %5d ", TimeNowOfSecond()-StartCycle)
	for _, v := range c {
		if modbusWork {
			fmt.Printf("%v ", v)
			chanIDS <- v
		}
	}
	fmt.Println(".")
	delete(counts, sc)
}
func TimeNowOfSecond() int {
	return time.Now().Hour()*3600 + time.Now().Minute()*60 + time.Now().Second()
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
			time.Sleep(time.Second)
			continue
		}
		client := modbus.NewClient(master)
		modbusWork = true

	work:
		for {
			ids := <-chanIDS
			if _, err := client.WriteSingleRegister(uint16(ids.ID), uint16(uint16(ids.Value))); err != nil {
				fmt.Printf("error modbus %s %s", con, err.Error())
				modbusWork = false
				break work
			}

		}
		master.Close()
	}

}
