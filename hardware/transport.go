package hardware

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kuda/modbus"
	"github.com/ruraomsk/kuda/setup"
)

type MasterTcp struct {
	master     modbus.TCPClientHandler
	client     modbus.Client
	hrInternal []uint16
	hr         []uint16
}

// type MasterRtu struct {
// 	master modbus.RTUClientHandler
// 	client modbus.Client
// 	hr     []uint16
// }

func (m *MasterTcp) readAllHR() error {
	if len(m.hrInternal) == 0 {
		return nil
	}
	for i := 0; i < len(m.hr); i++ {
		if m.hr[i] != m.hrInternal[i] {
			if _, err := m.client.WriteSingleRegister(uint16(i), m.hr[i]); err != nil {
				return err
			}
		}
	}
	size := uint16(100)
	ref := uint16(0)
	for count := len(m.hrInternal); count > 0; count -= int(size) {
		len := count
		if count > int(size) {
			len = int(size)
		}
		buff, err := m.client.ReadHoldingRegisters(ref, uint16(len))
		if err != nil {
			logger.Error.Printf("read hr %d %d %s", ref, len, err.Error())
			return err
		}
		pos := ref
		left := 0
		for i := 0; i < len; i++ {
			m.hrInternal[pos] = (uint16(buff[left]) << 8) | uint16(buff[left+1])
			pos++
			left += 2
		}
		ref += size
	}
	for i := 0; i < len(m.hr); i++ {
		m.hr[i] = m.hrInternal[i]
	}
	return nil
}

func (s *ModuleCPU) setMasterTCP() error {
	m := &MasterTcp{hrInternal: make([]uint16, s.size), hr: make([]uint16, s.size)} //549)}
	s.connect = "127.0.0.1:502"
	m.master = *modbus.NewTCPClientHandler(s.connect)
	m.master.SlaveId = byte(s.moduleSlaveID)
	m.master.Timeout = time.Second
	m.master.IdleTimeout = time.Minute
	if err := m.master.Connect(); err != nil {
		return fmt.Errorf("error modbus %s", err.Error())
	}
	m.client = modbus.NewClient(&m.master)
	s.masterTCP = m
	err := m.readAllHR()
	return err
}
func (s *ModuleCPU) loopTCP() {

	loop := time.NewTicker(time.Duration(setup.Set.Hardware.Step) * time.Millisecond)
	for {
	internal:
		for {
			select {
			case <-loop.C:
				if err := s.masterTCP.readAllHR(); err != nil {
					logger.Error.Printf("cpu %s", err.Error())
					s.work = false
					break internal
				} else {
					s.work = true
				}
			case wr := <-s.writer:
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
			s.masterTCP.master = *modbus.NewTCPClientHandler(s.connect)
			if err := s.masterTCP.master.Connect(); err != nil {
				logger.Error.Printf("error modbus %s", err.Error())
				time.Sleep(time.Second)
				continue
			}
			s.masterTCP.client = modbus.NewClient(&s.masterTCP.master)
			if err := s.masterTCP.readAllHR(); err != nil {
				logger.Error.Printf("error modbus %s", err.Error())
				time.Sleep(time.Second)
				continue
			}
			s.work = true
			break

		}
	}
}
