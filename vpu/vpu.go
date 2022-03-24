package vpu

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/setup"
	"github.com/ruraomsk/kuda/tech"
	"github.com/ruraomsk/kuda/tech/bin"
)

/*
Поддерживает обмен с ВПУ и контроллером

*/
var oldPhase int
var comPhase chan bin.PhaseCommand
var mutex sync.Mutex
var status SendToVPU

type SendToVPU struct {
	Phase int   `json:"phase"` //Текущая фаза
	Level int   `json:"level"` //Уровень управления 0-РУ 1-РПУ 2-ВПУ 3-ДУ 4-КУ
	Time  int   `json:"time"`  //Текущее время устройства сек от начала суток
	Naps  []Nap `json:"naps"`  //Состояние Направлений
}
type Nap struct {
	Number int  `json:"num"`
	Value  bool `json:"value"` //0 выключен 1 включен
}

func closeWorker(socket net.Conn) {
	comPhase <- bin.PhaseCommand{Level: 2, Phase: -1}
	logger.Info.Printf("User VPU %s closed", socket.RemoteAddr().String())
	socket.Close()
}
func workerVpu(socket net.Conn) {
	defer closeWorker(socket)
	reader := bufio.NewReader(socket)
	writer := bufio.NewWriter(socket)
	for {
		sphase, err := reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("Vpu user %s %s", socket.RemoteAddr().String(), err.Error())
			return
		}
		sphase = strings.Replace(sphase, "\n", "", -1)
		phase, err := strconv.Atoi(sphase)
		if err != nil {
			logger.Error.Printf("Vpu user %s send '%s'", socket.RemoteAddr().String(), sphase)
			phase = 0
		}
		if oldPhase != phase {
			if phase == 0 {
				comPhase <- bin.PhaseCommand{Level: 2, Phase: -1}
			} else {
				comPhase <- bin.PhaseCommand{Level: 2, Phase: phase}
			}
			oldPhase = phase
			fmt.Printf("vpu new cmd  %d\n", phase)
		}
		mutex.Lock()
		status.Time = bin.TimeNowOfSecond()
		status.Naps = make([]Nap, 0)
		for n, v := range tech.CmkNow.Naps {
			status.Naps = append(status.Naps, Nap{Number: n, Value: v})
		}
		t := status
		mutex.Unlock()
		buffer, err := json.Marshal(&t)
		if err != nil {
			logger.Error.Printf("json %s", err.Error())
			buffer = make([]byte, 0)
		}
		writer.WriteString(string(buffer) + "\n")
		writer.Flush()
	}
}
func ListenExternalVpu() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%d", setup.Set.Vpu.SPort))
	if err != nil {
		logger.Error.Printf("Open VPU port %s", err.Error())
		return
	}
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Accept %s", err.Error())
			continue
		}
		logger.Info.Printf("new user VPU %s", socket.RemoteAddr().String())
		go workerVpu(socket)
	}
}
func VpuExch(cm chan bin.ResponcePhase, cf chan bin.PhaseCommand) {
	comPhase = cf
	ListenExternalVpu()
	oldPhase = 0
	for {
		resp := <-cm
		mutex.Lock()
		status.Phase = resp.Phase
		status.Level = resp.Level
		mutex.Unlock()
	}

}
func StarterVPU(cm *tech.Commander, cf chan bin.PhaseCommand) {
	st := bin.ResponcePhase{}
	cmm := make(chan bin.ResponcePhase)
	go VpuExch(cmm, cf)
	for {
		for {
			resp := <-cm.Responce
			if !reflect.DeepEqual(&st, &resp) {
				cmm <- resp
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
