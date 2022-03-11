package bin

import "fmt"

type Command struct {
	Tir   int //Номер тиристора
	Value int //Значение 0 выкл 1 включить
}
type PromMake struct {
	Ticks map[int][]Command
}

//InitCMK иницилизирует и проверяет на полноту
func (c *CMK) InitCMK() error {
	//Проверяем и считаем направления
	c.Naps = make(map[int]bool)
	c.PromMake.new()
	if c.StepPromtact == 0 {
		c.StepPromtact = 100
	}
	for _, v := range c.TirToNaps {
		v.setDefault()
		c.Naps[v.Number] = false
	}
	//Проверяем фазы на полноту
	for _, v := range c.NtoPhases {
		for _, w := range v.Naps {
			_, is := c.Naps[w]
			if !is {
				return fmt.Errorf("into phase %d is nap % but not define", v.NumPhase, w)
			}
		}
	}
	//Проверяем промтакты и кешируем их
	c.PrMapBase = make(map[int]PromTakt)
	c.PrMapUn = make(map[int]PromTakt)
	for _, v := range c.PromTaktBases {
		_, is := c.Naps[v.Nap]
		if !is {
			return fmt.Errorf("is nap % but not define protak base", v.Nap)
		}
		c.PrMapBase[v.Nap] = v
	}
	for _, v := range c.PromTakt {
		_, is := c.Naps[v.Nap]
		if !is {
			return fmt.Errorf("is nap % but not define protak univer", v.Nap)
		}
		c.PrMapUn[v.Nap] = v
	}
	return nil
}
func (c *CMK) GetTMin(phase int) int {
	for _, v := range c.TminToPhases {
		if v.NumPhase == phase {
			return v.Tmin
		}
	}
	return 20
}
func (p *PromMake) new() {
	p.Ticks = make(map[int][]Command)
}
func (p *PromMake) add(t int, cmd Command) {
	_, is := p.Ticks[t]
	if !is {
		p.Ticks[t] = make([]Command, 0)

	}
	p.Ticks[t] = append(p.Ticks[t], cmd)
}
func (p *PromMake) GetMaxTime() int {
	tmax := 0
	for t := range p.Ticks {
		if t > tmax {
			tmax = t
		}
	}
	return tmax
}
func (p *PromMake) GetCommads(t int) []Command {
	c, is := p.Ticks[t]
	if !is {
		return make([]Command, 0)
	}
	return c
}

func (t *TirToNap) setDefault() {
	if t.Type == 0 {
		if t.Green != 0 && t.Yellow != 0 && len(t.Reds) != 0 {
			//Это транспортное направление
			t.Type = 1
			return
		}
		if t.Green != 0 && t.Yellow == 0 && len(t.Reds) != 0 {
			//Это пешеходный
			t.Type = 3
			return
		}
		if t.Green != 0 && t.Yellow == 0 && len(t.Reds) == 0 {
			//Это поворотный
			t.Type = 2
			return
		}
	}
}
func (t *TirToNap) getYellowOn() []Command {
	res := make([]Command, 0)
	for _, v := range t.Reds {
		res = append(res, Command{Tir: v, Value: 0})
	}
	if t.Green > 0 {
		res = append(res, Command{Tir: t.Green, Value: 0})
	}
	if t.Yellow > 0 {
		res = append(res, Command{Tir: t.Yellow, Value: 1})
	}
	return res
}

func (t *TirToNap) getYellowOff() []Command {
	res := make([]Command, 0)
	for _, v := range t.Reds {
		res = append(res, Command{Tir: v, Value: 0})
	}
	if t.Green > 0 {
		res = append(res, Command{Tir: t.Green, Value: 0})
	}
	if t.Yellow > 0 {
		res = append(res, Command{Tir: t.Yellow, Value: 0})
	}
	return res
}

//SetAllYellowOn - возвращает команды которые отыгрывают первую половину желтого мигания
func (c *CMK) SetAllYellowOn() []Command {
	res := make([]Command, 0)
	for _, v := range c.TirToNaps {
		res = append(res, v.getYellowOn()...)
	}
	return res
}

//SetAllYellowOff - возвращает команды которые отыгрывают вторую половину желтого мигания
func (c *CMK) SetAllYellowOff() []Command {
	res := make([]Command, 0)
	for _, v := range c.TirToNaps {
		res = append(res, v.getYellowOff()...)
	}
	return res
}

//SetAllRed - возращает команды для включения режима все Красные
func (c *CMK) SetAllRed() []Command {
	res := make([]Command, 0)
	for _, v := range c.TirToNaps {
		res = append(res, v.getRedOn()...)
	}
	return res
}
func (c *CMK) SetAllNapsStop() {
	for n := range c.Naps {
		c.Naps[n] = false
	}
}

//SetAllOff - возвращает команды выключения всех тиристоров
func (c *CMK) SetAllOff() []Command {
	res := make([]Command, 0)
	for _, v := range c.TirToNaps {
		res = append(res, v.getAllOff()...)
	}
	return res
}

func (t *TirToNap) getRedOn() []Command {
	res := make([]Command, 0)
	for _, v := range t.Reds {
		res = append(res, Command{Tir: v, Value: 1})
	}
	if t.Green > 0 {
		res = append(res, Command{Tir: t.Green, Value: 0})
	}
	if t.Yellow > 0 {
		res = append(res, Command{Tir: t.Yellow, Value: 0})
	}
	return res
}
func (t *TirToNap) getAllOff() []Command {
	res := make([]Command, 0)
	for _, v := range t.Reds {
		res = append(res, Command{Tir: v, Value: 0})
	}
	if t.Green > 0 {
		res = append(res, Command{Tir: t.Green, Value: 0})
	}
	if t.Yellow > 0 {
		res = append(res, Command{Tir: t.Yellow, Value: 0})
	}
	return res
}
