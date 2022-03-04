package tech

import (
	"fmt"

	"github.com/ruraomsk/kuda/tech/bin"
)

var (
	ps plans
)

func (p *plans) addValue(t int, v value) {
	ps, is := p.maps[t]
	if !is {
		m := new(plan)
		m.values = make([]value, 0)
		p.maps[t] = m
		ps = p.maps[t]
	}
	ps.values = append(ps.values, v)
	p.maps[t] = ps
}

type plans struct {
	maps map[int]*plan
}

type plan struct {
	values []value
}
type value struct {
	chanel int //Номер выхода
	val    int //Новое значение
}

func makePlanPhase(phase, tphase int) (start plan, pls map[int]*plan, err error) {
	start = plan{values: make([]value, 0)}
	pls = make(map[int]*plan)
	ps.maps = make(map[int]*plan)
	err = nil
	// Первым делом создаем стартовое состяние зеленый на работающих красный на стоящих
	// Составляем полный список направлений
	tirtonaps := make(map[int]bin.TirToNap)
	tirdone := make(map[int]bool)
	for _, v := range cmk.TirToNaps {
		tirtonaps[v.Number] = v
		tirdone[v.Number] = false
	}
	//Пробегаем по фазам работающим
	found := false
	for _, v := range cmk.NtoPhases {
		if v.NumPhase == phase {
			for _, n := range v.Naps {
				t := tirtonaps[n]
				start.values = append(start.values, value{t.Green, 1})
				start.values = append(start.values, value{t.Yellow, 0})
				for _, r := range t.Reds {
					start.values = append(start.values, value{r, 0})
				}
				tirdone[n] = true
			}
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("phase %d not found", phase)
		return
	}
	for n, v := range tirdone {
		if !v {
			t := tirtonaps[n]
			start.values = append(start.values, value{t.Green, 0})
			start.values = append(start.values, value{t.Yellow, 0})
			for _, r := range t.Reds {
				start.values = append(start.values, value{r, 1})
			}

		}
	}
	// Начинаем делать промтакт по всем
	tirtonaps = make(map[int]bin.TirToNap)
	tirdone = make(map[int]bool)
	proms := make(map[int]bin.PromTakt)
	for _, v := range cmk.TirToNaps {
		tirtonaps[v.Number] = v
		tirdone[v.Number] = false
	}
	for _, v := range cmk.PromTakt {
		proms[v.Nap] = v
	}
	//Переход по открытым направлениям
	for _, v := range cmk.NtoPhases {
		if v.NumPhase == phase {
			for _, n := range v.Naps {
				t := tirtonaps[n]
				pr := proms[n]
				if pr.GreenDop > pr.GreenBlink {
					ps.addValue(pr.GreenDop*100, value{t.Green, 1})
				}
				if pr.GreenBlink > pr.Yellow {
					for i := pr.GreenBlink; i < pr.Yellow; i++ {
						ps.addValue(i*100, value{t.Green, 0})
						ps.addValue((i*100)+500, value{t.Green, 1})
					}

				}
				if pr.Yellow > 0 {
					//Просто ставим время начала желтого
					ps.addValue(pr.Yellow*100, value{t.Green, 0})
					ps.addValue(pr.Yellow*100, value{t.Yellow, 1})
				}
				tirdone[n] = true
			}
			break
		}
	}

	// Переход по закрытым направлениям
	for n, v := range tirdone {
		if !v {
			t := tirtonaps[n]
			pr := proms[n]
			if pr.Red > pr.RedYellow {
				for _, r := range t.Reds {
					ps.addValue(pr.Red*100, value{r, 1})
				}

			}
			if pr.RedYellow > 0 {
				ps.addValue(pr.RedYellow*100, value{t.Yellow, 1})
				for _, r := range t.Reds {
					ps.addValue(pr.RedYellow*100, value{r, 1})
				}
			}
		}
	}

	return
}
