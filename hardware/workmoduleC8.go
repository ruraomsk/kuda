package hardware

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/modbus"
	"github.com/ruraomsk/kuda/setup"
)

func (s *ModuleC8) setMasterTCP() error {
	m := &MasterTcp{hrInternal: make([]uint16, s.size), hr: make([]uint16, s.size)} //549)}
	for i := 0; i < len(m.hr); i++ {
		m.hr[i] = 0
		m.hrInternal[i] = 0
	}
	m.master = *modbus.NewTCPClientHandler(s.connect)
	m.master.SlaveId = byte(s.moduleSlaveID)
	m.master.Timeout = time.Second
	m.master.IdleTimeout = time.Minute
	s.writer = make(chan writeHR)
	if err := m.master.Connect(); err != nil {
		return fmt.Errorf("error modbus %s %s", s.connect, err.Error())
	}
	m.client = modbus.NewClient(&m.master)
	s.masterTCP = m
	err := m.readAllHR()
	return err
}
func (s *ModuleC8) loopTCP() {
	// logger.Info.Printf("start modbus loop %s", s.connect)
	loop := time.NewTicker(time.Duration(setup.Set.Hardware.Step) * time.Millisecond)
	for {
	internal:
		for {
			select {
			case <-loop.C:
				if err := s.masterTCP.readAllHR(); err != nil {
					s.work = false
					break internal
				} else {
					s.work = true
				}
			case wr := <-s.writer:
				// logger.Debug.Printf("%d %v", s.moduleNumber, wr)
				//Пришла команда на запись если поле bit <0 то это просто слово
				if wr.pos.b < 0 {
					s.masterTCP.hr[wr.pos.w] = uint16(wr.value)
					s.work = true
				} else {
					r := s.masterTCP.hr[wr.pos.w]
					c := 1
					for i := 0; i < wr.pos.b; i++ {
						c = c << 1
					}
					if wr.value > 0 {
						r = r | uint16(c)
					} else {
						r = r & (^uint16(c))
					}
					s.masterTCP.hr[wr.pos.w] = r
					s.work = true
				}
			}
		}
		s.masterTCP.master.Close()
		time.Sleep(time.Second)

		for {
			if err := s.masterTCP.master.Connect(); err != nil {
				logger.Error.Printf("error modbus %s %s", s.connect, err.Error())
				time.Sleep(10 * time.Second)
				continue
			}
			s.masterTCP.client = modbus.NewClient(&s.masterTCP.master)
			if err := s.masterTCP.readAllHR(); err != nil {
				logger.Error.Printf("error modbus %s %s", s.connect, err.Error())
				time.Sleep(10 * time.Second)
				continue
			}
			s.work = true
			break
		}

	}

}
