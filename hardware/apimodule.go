package hardware

import "github.com/ruraomsk/ag-server/logger"

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
	// logger.Debug.Printf("%v %d", v.value, value)
	m.writer <- writeHR{pos: v.value, value: value}
}
