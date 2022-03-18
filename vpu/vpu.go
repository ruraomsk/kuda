package vpu

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/modbus"
	"github.com/ruraomsk/kuda/setup"
	"github.com/ruraomsk/kuda/tech"
	"github.com/ruraomsk/kuda/tech/bin"
)

/*
Поддерживает обмен с ВПУ и контроллером

*/
var modbusWork = false

func modbusExch(cm chan bin.ResponcePhase, cf chan bin.PhaseCommand) {
	con := fmt.Sprintf("%s:%d", setup.Set.Vpu.Connect, setup.Set.Vpu.SPort)
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
		ticks := time.NewTicker(time.Duration(setup.Set.Vpu.Step) * time.Millisecond)
		oldPhase := uint16(0)
		modbusWork = true
	work:
		for {
			select {
			case <-ticks.C:
				buff, err := client.ReadHoldingRegisters(0, 1)
				if err != nil {
					logger.Error.Printf("read hr %s %s", con, err.Error())
					modbusWork = false
					break work
				}
				newPhase := (uint16(buff[0]) << 8) | uint16(buff[1])
				if oldPhase != newPhase {
					if newPhase == 0 {
						cf <- bin.PhaseCommand{Level: 2, Phase: -1}
					} else {
						cf <- bin.PhaseCommand{Level: 2, Phase: int(newPhase)}
					}
					oldPhase = newPhase
					fmt.Printf("vpu new cmd  %d\n", newPhase)
				}

			case resp := <-cm:
				if _, err := client.WriteSingleRegister(1, uint16(resp.Phase)); err != nil {
					fmt.Printf("error modbus %s %s", con, err.Error())
					modbusWork = false
					break work
				}
				if _, err := client.WriteSingleRegister(2, uint16(resp.Level)); err != nil {
					modbusWork = false
					fmt.Printf("error modbus %s %s", con, err.Error())
					modbusWork = false
					break work
				}
				// fmt.Printf("vpu: %v\n", resp)
			}

		}
		master.Close()
	}

}
func StarterVPU(cm *tech.Commander, cf chan bin.PhaseCommand) {
	st := bin.ResponcePhase{}
	cmm := make(chan bin.ResponcePhase)
	go modbusExch(cmm, cf)
	for {
		for {
			resp := <-cm.Responce
			if !reflect.DeepEqual(&st, &resp) {
				if modbusWork {
					cmm <- resp
				}
				// fmt.Printf("vpu: %v\n", resp)
				st = resp
			}

		}
	}
}
func StarterDU(cm *tech.Commander, cf chan bin.PhaseCommand) {
	st := bin.ResponcePhase{}
	for {
		resp := <-cm.Responce
		if !reflect.DeepEqual(&st, &resp) {
			// fmt.Printf("KU : %v\n", resp)
			st = resp
		}

	}
}
func StarterKU(cm *tech.Commander, cf chan bin.PhaseCommand) {
	st := bin.ResponcePhase{}
	for {
		resp := <-cm.Responce
		if !reflect.DeepEqual(&st, &resp) {
			// fmt.Printf("KU : %v\n", resp)
			st = resp
		}

	}
}
