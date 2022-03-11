package hardware

import (
	"github.com/ruraomsk/ag-server/logger"
)

func (c *ModuleCPU) IsWork() bool {
	return c.work
}
func (c *ModuleCPU) GetDI(nomer int) bool {
	v, is := c.di[nomer]
	if !is {
		logger.Error.Printf("di getvalue bad nomer %d", nomer)
		return false
	}
	return v.getValue(c.masterTCP.hr)
}
func (c *ModuleCPU) SetDO(nomer int, value int) {
	v, is := c.do[nomer]
	if !is {
		logger.Error.Printf("do setvalue bad nomer %d", nomer)
		return
	}
	if c.masterTCP.hr[v.state] != 201 {
		logger.Error.Printf("do setvalue nomer %d not DO %d", nomer, c.masterTCP.hr[v.state])
		return
	}
	c.writer <- writeHR{pos: v.value, value: value}
}
func IsWorkC8(number int) bool {
	m, is := MapC8[number]
	if !is {
		logger.Error.Printf("c8 number bad  %d", number)
		return false

	}
	return m.work
}
func C8GetValue(number, chanel int) bool {
	m, is := MapC8[number]
	if !is {
		logger.Error.Printf("c8 number bad  %d", number)
		return false

	}
	v, is := m.с8[chanel]
	if !is {
		logger.Error.Printf("c8 number %d bad chanel %d", number, chanel)
		return false
	}
	return v.getValue(m.masterTCP.hr)
}
func C8SetValue(number, chanel, value int) {
	m, is := MapC8[number]
	if !is {
		logger.Error.Printf("c8 number bad  %d", number)
		return

	}
	v, is := m.с8[chanel]
	if !is {
		logger.Error.Printf("c8 number %d bad chanel %d", number, chanel)
		return
	}
	m.writer <- writeHR{pos: v.value, value: value}
}
func C8GetOut(chanel int) bool {
	number := ((chanel - 1) / 8) + 2
	m, is := MapC8[number]
	if !is {
		logger.Error.Printf("c8 number bad  %d", number)
		return false

	}
	ch := chanel % 8
	if ch == 0 {
		ch = 8
	}
	v, is := m.с8[ch]
	if !is {
		logger.Error.Printf("c8 number %d bad chanel %d", number, ch)
		return false
	}
	return v.getValue(m.masterTCP.hr)
}
func C8SetOut(chanel, value int) {
	number := ((chanel - 1) / 8) + 2
	m, is := MapC8[number]
	if !is {
		logger.Error.Printf("c8 number bad  %d", number)
		return

	}
	ch := chanel % 8
	if ch == 0 {
		ch = 8
	}
	v, is := m.с8[ch]
	if !is {
		logger.Error.Printf("c8 number %d bad chanel %d", number, chanel%8)
		return
	}
	m.writer <- writeHR{pos: v.value, value: value}
}
