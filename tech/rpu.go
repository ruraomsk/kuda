package tech

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/hardware"
	"github.com/ruraomsk/kuda/tech/bin"
)

var cmk *bin.CMK
var commands chan bin.PhaseCommand
var responce chan bin.ResponcePhase
var err error

func WorkRPU(c *bin.CMK) {
	cmk = c
	for !hardware.AllReady() {
		fmt.Println("wait...")
		time.Sleep(time.Second)
	}
	commands = make(chan bin.PhaseCommand)
	responce, err = bin.StartMechanics(cmk, commands)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	step := 0
	phaseNow := 0
	ctrlCycle := time.NewTimer(1 * time.Hour)
	ctrlPhase := time.NewTimer(100 * time.Hour)
	for {
		select {
		case <-ctrlCycle.C:
			step = 0
			ctrlCycle = time.NewTimer(time.Duration(cmk.RPUs[0].Tcycle) * time.Second)
			ctrlPhase = time.NewTimer(time.Duration(cmk.RPUs[0].Phases[0].Time) * time.Second)
			phaseNow = cmk.RPUs[0].Phases[0].Phase
			commands <- bin.PhaseCommand{Phase: phaseNow, PromTakt: true}
		case <-ctrlPhase.C:
			step++
			ctrlPhase = time.NewTimer(time.Duration(cmk.RPUs[0].Phases[step].Time) * time.Second)
			phaseOld := phaseNow
			phaseNow = cmk.RPUs[0].Phases[step].Phase
			commands <- bin.PhaseCommand{Phase: phaseNow, PromTakt: cmk.GetBaseOrUniver(phaseOld, phaseNow)}
		case resp := <-responce:
			if resp.OsStop {
				fmt.Println("RU OS Stop")
				ctrlCycle.Stop()
				ctrlPhase.Stop()
				continue
			}
			if resp.YellowBlink {
				fmt.Println("RU Yellow Blink")
				ctrlCycle.Stop()
				ctrlPhase.Stop()
				continue
			}
			if resp.Ready && resp.Phase == 12 {
				fmt.Println("All red ready work")
				ctrlCycle = time.NewTimer(1 * time.Millisecond)
				continue
			}
			// if !resp.Ready {
			fmt.Printf("responce %v\n", resp)
			// }

		}

	}

}
