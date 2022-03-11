package bin

import (
	"fmt"
	"math"
	"time"

	"github.com/ruraomsk/kuda/hardware"
	"github.com/ruraomsk/kuda/setup"
)

//StartMechanics основной исполнитель

type PhaseCommand struct {
	Phase    int
	PromTakt bool //Истина если используется базовый промтакт
}
type ResponcePhase struct {
	Phase       int //Номер исполняемой фазы
	Ready       bool
	OsStop      bool // Истина если включен режим РУ ОС
	YellowBlink bool // Истина если включен режим РУ ЖМ
}

var (
	cmk            *CMK
	phaseNow       int
	timeCount      int
	timePhase      int
	YellowOnOrOff  bool
	inPhaseCommand chan PhaseCommand
	responce       chan ResponcePhase
)

func StartMechanics(c *CMK, phase chan PhaseCommand) (chan ResponcePhase, error) {
	responce = make(chan ResponcePhase, 100)
	cmk = c
	if err := cmk.InitCMK(); err != nil {
		return responce, err
	}
	inPhaseCommand = phase
	go mainWork()
	return responce, nil
}

func mainWork() {
	for !hardware.AllReady() {
		fmt.Println("wait...")
		time.Sleep(time.Second)
	}
hardosStop:
	hardOsWork()
hardallBlink:
	hardAllYellowBlink()
	writeValues(cmk.SetAllRed())
	time.Sleep(time.Duration(time.Duration(setup.Set.Hardware.LongKK) * time.Second))
	responce <- ResponcePhase{Phase: 12, Ready: true}
	//Все направления выключены
	cmk.SetAllNapsStop()
	timeCount = -cmk.StepPromtact
	controlSwitches := time.NewTicker(time.Duration(cmk.StepPromtact) * time.Millisecond)
	allHalfSecond := time.NewTicker(500 * time.Millisecond)
	oneSecond := time.NewTicker(time.Second)
	controlPromtakt := time.NewTicker(time.Duration(cmk.StepPromtact) * time.Millisecond)
	timePhase = -1
	once := false
	for {
		select {
		case <-oneSecond.C:
			if timePhase >= 0 {
				timePhase++
				if timePhase > (math.MaxInt - 10) {
					timePhase = 0
				}
			}
		case <-controlPromtakt.C:
			if timeCount >= 0 {
				if !once {
					responce <- ResponcePhase{Phase: 0, Ready: true}
					once = true
				}
				p := cmk.PromMake.GetCommads(timeCount)
				if len(p) != 0 {
					fmt.Printf("time %4d : ", timeCount)
					for _, v := range p {
						hardware.C8SetOut(v.Tir, v.Value)
						fmt.Printf("%v", v)
					}
					fmt.Println(".")
				}
				timePhase = -1
			}
			if timeCount == 0 {
				responce <- ResponcePhase{Phase: phaseNow, Ready: true}
				timePhase = 0
			}
			timeCount -= cmk.StepPromtact
			if timeCount < 0 {
				timeCount = -cmk.StepPromtact
			}

		case <-allHalfSecond.C:
			switch phaseNow {
			case 10:
				if !once {
					responce <- ResponcePhase{Phase: 10, Ready: true}
					once = true
				}
				if YellowOnOrOff {
					writeValues(cmk.SetAllYellowOn())
				} else {
					writeValues(cmk.SetAllYellowOff())
				}
				YellowOnOrOff = !YellowOnOrOff
			case 11:
				if !once {
					responce <- ResponcePhase{Phase: 11, Ready: true}
					once = true
					writeValues(cmk.SetAllOff())
				}
			case 12:
				if !once {
					writeValues(cmk.SetAllRed())
					time.Sleep(time.Duration(time.Duration(setup.Set.Hardware.LongKK) * time.Second))
					responce <- ResponcePhase{Phase: 12, Ready: true}
					once = true
				}
			}
		case <-controlSwitches.C:
			if hardware.Cpu.GetDI(setup.Set.Hardware.PinOS) {
				controlSwitches.Stop()
				goto hardosStop
			}
			if hardware.Cpu.GetDI(setup.Set.Hardware.PinYB) {
				controlSwitches.Stop()
				goto hardallBlink
			}
		case cmd := <-inPhaseCommand:
			fmt.Printf("command %v\n", cmd)
			switch cmd.Phase {
			case 10:
				YellowOnOrOff = true
				phaseNow = 10
				once = false
			case 11:
				phaseNow = 11
				once = false
			case 12:
				phaseNow = 12
				once = false
				//Все направления выключены
				cmk.SetAllNapsStop()

			default:
				if phaseNow == 10 || phaseNow == 11 {
					writeValues(cmk.SetAllRed())
					time.Sleep(time.Duration(time.Duration(setup.Set.Hardware.LongKK) * time.Second))
					responce <- ResponcePhase{Phase: 12, Ready: true}
					//Все направления выключены
					cmk.SetAllNapsStop()
					continue
				}
				if timeCount > 0 {
					//Еще идет промтакт посылаем отказ
					responce <- ResponcePhase{Phase: 0, Ready: false}
					continue
				}
				if timePhase <= cmk.GetTMin(phaseNow) {
					//Еще не выбран Тмин текущей фазы
					responce <- ResponcePhase{Phase: phaseNow, Ready: false}
					continue
				}
				err := cmk.GetPromtackt(cmd.Phase, cmd.PromTakt)
				if err != nil {
					responce <- ResponcePhase{Phase: cmd.Phase, Ready: false}
					continue
				}
				timeCount = cmk.PromMake.GetMaxTime()
				once = false
				phaseNow = cmd.Phase
				timePhase = 0
			}

		}
	}
}
func writeValues(cmds []Command) {
	for _, v := range cmds {
		hardware.C8SetOut(v.Tir, v.Value)
	}
}
func zeroOn() {
	writeValues(cmk.SetAllOff())
}
func hardOsWork() bool {
	result := false
	once := false
	for hardware.Cpu.GetDI(setup.Set.Hardware.PinOS) {
		if !once {
			once = true
			responce <- ResponcePhase{OsStop: true}
		}
		result = true
		zeroOn()
		time.Sleep(100 * time.Millisecond)
	}
	return result
}

func hardAllYellowBlink() {
	once := false
	for hardware.Cpu.GetDI(setup.Set.Hardware.PinYB) {
		if !once {
			once = true
			responce <- ResponcePhase{YellowBlink: true}
		}
		if hardOsWork() {
			continue
		}
		writeValues(cmk.SetAllYellowOn())
		time.Sleep(500 * time.Millisecond)
		if hardOsWork() {
			continue
		}
		writeValues(cmk.SetAllYellowOff())
		time.Sleep(500 * time.Millisecond)
		if hardOsWork() {
			continue
		}
	}
	zeroOn()
}
