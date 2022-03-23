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

type Commander struct {
	Level       int //Источник команд 0-РУ 1-РПУ 2-ВПУ 3-ДУ 4-КУ
	Responce    chan bin.ResponcePhase
	CommandWait bin.PhaseCommand
}

var levelNow int
var coms map[int]*Commander

func CreateCommander(level int) *Commander {
	return &Commander{Level: level, Responce: make(chan bin.ResponcePhase), CommandWait: bin.PhaseCommand{Phase: -1}}
}
func WorkRPU(c *bin.CMK, cs []*Commander, CommandFlow chan bin.PhaseCommand) {
	cmk = c
	coms = make(map[int]*Commander)
	coms[0] = &Commander{Level: 0, Responce: make(chan bin.ResponcePhase), CommandWait: bin.PhaseCommand{Phase: 1, PromTakt: true, LongTime: cmk.RPUs[0].Phases[0].Time}}
	for _, v := range cs {
		coms[v.Level] = v
	}
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
	resp := <-responce
	for !(resp.Ready && resp.Phase == 12) {
		resp = <-responce
	}

	levelNow = 1
	step := 0
	phaseNow := 12

	ctrlCycle := time.NewTimer(100 * time.Hour)
	ctrlPhase := time.NewTimer(100 * time.Hour)
	seekLevel := time.NewTicker(100 * time.Millisecond)

	for {
		select {
		case <-seekLevel.C:
			if phaseNow != 9 {
				lev := -1
				for l := 2; l < 5; l++ {
					lc, is := coms[l]
					if !is {
						continue
					}
					if lc.CommandWait.Phase != -1 {
						lev = l
						break
					}
				}
				if lev < 0 {
					lev = 0
				}
				// fmt.Printf("select level %d\n", lev)
				send := true
				if lev != 0 {
					ctrlCycle.Stop()
					ctrlPhase.Stop()
				} else {
					if lev == 0 && levelNow != 0 {
						coms[0].CommandWait = bin.PhaseCommand{Level: 0, Phase: 1, PromTakt: true, LongTime: cmk.RPUs[0].Phases[0].Time}
						fmt.Printf("need start cycle %5d\n", bin.TimeNowOfSecond())
						cmk.MakeCycleGramm(0)
						commands <- coms[0].CommandWait
						ctrlCycle = time.NewTimer(time.Duration(cmk.RPUs[0].Tcycle) * time.Second)
						ctrlPhase = time.NewTimer(time.Duration(cmk.RPUs[0].Phases[0].Time) * time.Second)
						send = false
					}
				}
				levelNow = lev
				if send {
					lc := coms[lev]
					commands <- bin.PhaseCommand{Level: levelNow, Phase: lc.CommandWait.Phase, PromTakt: cmk.GetBaseOrUniver(phaseNow, lc.CommandWait.Phase), LongTime: lc.CommandWait.LongTime}
				}
			}

		case cmd := <-CommandFlow:
			cl, is := coms[cmd.Level]
			if !is {
				logger.Error.Printf("Not found level %d", cmd.Level)
				continue
			}
			// fmt.Printf("new cmd %v\n", cmd)
			if cmd.Phase < 0 {
				cl.CommandWait = cmd
				continue
			}
			if cmd.Phase == 10 || cmd.Phase == 11 || cmd.Phase == 12 {
				cl.CommandWait = cmd
				continue
			}
			if cmk.IsPhase(cmd.Phase) {
				cl.CommandWait = cmd
			}
		case <-ctrlCycle.C:
			if levelNow == 0 {
				step = 0
				ctrlCycle = time.NewTimer(time.Duration(cmk.RPUs[0].Tcycle) * time.Second)
				ctrlPhase = time.NewTimer(time.Duration(cmk.RPUs[0].Phases[0].Time) * time.Second)
				if phaseNow != cmk.RPUs[0].Phases[0].Phase {
					phaseNow = cmk.RPUs[0].Phases[0].Phase
					fmt.Printf("need start cycle %5d\n", bin.TimeNowOfSecond())
					cmk.MakeCycleGramm(0)
					CommandFlow <- bin.PhaseCommand{Level: 0, Phase: phaseNow, PromTakt: true, LongTime: cmk.RPUs[0].Phases[0].Time}
				}
			}
		case <-ctrlPhase.C:
			if levelNow == 0 {
				step++
				if step >= len(cmk.RPUs[0].Phases) {
					continue
				}
				fmt.Printf("%5d end phase %d\n", bin.TimeNowOfSecond()-bin.StartCycle, phaseNow)
				ctrlPhase = time.NewTimer(time.Duration(cmk.RPUs[0].Phases[step].Time) * time.Second)
				phaseOld := phaseNow
				phaseNow = cmk.RPUs[0].Phases[step].Phase
				CommandFlow <- bin.PhaseCommand{Phase: phaseNow, PromTakt: cmk.GetBaseOrUniver(phaseOld, phaseNow), LongTime: cmk.RPUs[0].Phases[step].Time}
				fmt.Printf("%5d start phase %d\n", bin.TimeNowOfSecond()-bin.StartCycle, phaseNow)
			}
		case resp := <-responce:
			if resp.OsStop {
				// fmt.Println("RU OS Stop")
				sendOsStop()
				levelNow = 1
				ctrlCycle.Stop()
				ctrlPhase.Stop()

			}
			if resp.YellowBlink {
				// fmt.Println("RU Yellow Blink")
				sendYellowBlink()
				ctrlCycle.Stop()
				ctrlPhase.Stop()
				levelNow = 1
			}
			if resp.Phase == 12 {
				// fmt.Println("All red ready work")
				sendAllRed(resp)
				levelNow = 1
				ctrlCycle.Stop()
				ctrlPhase.Stop()
			}
			sendResponce(resp)
			if resp.Phase != 9 {
				phaseNow = resp.Phase
			}
		}

	}

}
func sendOsStop() {
	for _, v := range coms {
		if v.Level == 0 {
			continue
		}
		v.Responce <- bin.ResponcePhase{Level: 1, Phase: 11, OsStop: true}
	}
}
func sendYellowBlink() {
	for _, v := range coms {
		if v.Level == 0 {
			continue
		}
		v.Responce <- bin.ResponcePhase{Level: 1, Phase: 10, YellowBlink: true}
	}
}
func sendAllRed(resp bin.ResponcePhase) {
	for _, v := range coms {
		if v.Level == 0 {
			continue
		}
		v.Responce <- resp
	}
}
func sendResponce(resp bin.ResponcePhase) {
	for _, v := range coms {
		if v.Level == 0 {
			continue
		}
		v.Responce <- resp
	}

}
