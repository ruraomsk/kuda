package bin

import "fmt"

//GetPromtackt вычисляет переходной промтакт на фазу phase.prom истина если используется базовый промтакт
//	naps текущее состояние направлений Истина если открыто. На выходе изменяем состояние направлений и выдаем
//  план перехода  промтакта который нужно отработать и ошибку
func (c *CMK) GetPromtackt(phase int, prom bool) error {
	//Вначале проверим есть ли такая фаза?
	found := false
	descPhase := c.NtoPhases[0]
	for _, v := range c.NtoPhases {
		if v.NumPhase == phase {
			descPhase = v
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("phase %d not found", phase)
	}
	pr := c.PrMapBase
	if !prom {
		pr = c.PrMapUn
	}
	//Строим состяние по окончанию перехода
	newNaps := c.Naps
	for v := range newNaps {
		newNaps[v] = false
	}
	for _, v := range descPhase.Naps {
		newNaps[v] = true
	}

	for _, v := range c.TirToNaps {
		res := v.makeProm(c.Naps[v.Number], newNaps[v.Number], pr[v.Number])
		for t, k := range res.Ticks {
			for _, l := range k {
				c.PromMake.add(t, l)
			}
		}
	}
	c.Naps = newNaps
	return nil
}
func (tm *TirToNap) makeProm(olds bool, news bool, pr PromTakt) PromMake {
	res := new(PromMake)
	res.new()
	if olds == news {
		//Ничего не изменилось
		return *res
	}
	switch tm.Type {
	case 1: //Транспортное направление
		if olds {
			//Нужно закрыть направление Зм Ж Кр
			if pr.GreenBlink > pr.Yellow {
				//Записываем зеленое мигание до желтого
				for i := pr.GreenBlink; i > pr.Yellow; i-- {
					res.add(i*1000, Command{Tir: tm.Green, Value: 0})
					res.add(i*1000-500, Command{Tir: tm.Green, Value: 1})
				}
			}
			if pr.Yellow > pr.Red {
				res.add(pr.Yellow*1000, Command{Tir: tm.Green, Value: 0})
				res.add(pr.Yellow*1000, Command{Tir: tm.Yellow, Value: 1})
			}
			if pr.Red > 0 {
				res.add(pr.Red*1000, Command{Tir: tm.Yellow, Value: 0})
				for _, v := range tm.Reds {
					res.add(pr.Red*1000, Command{Tir: v, Value: 1})
				}
			}
			//Финалочка
			res.add(0, Command{Tir: tm.Yellow, Value: 0})
			res.add(0, Command{Tir: tm.Green, Value: 0})
			for _, v := range tm.Reds {
				res.add(0, Command{Tir: v, Value: 1})
			}

		} else {
			//Нужно открыть направление Кр КрЖ З
			if pr.RedYellow != 0 && pr.RedYellow > pr.GreenDop {
				res.add(pr.RedYellow*1000, Command{Tir: tm.Yellow, Value: 1})
				for _, v := range tm.Reds {
					res.add(pr.Red*1000, Command{Tir: v, Value: 1})
				}
			}
			if pr.GreenDop > 0 {
				res.add(pr.GreenDop*1000, Command{Tir: tm.Yellow, Value: 0})
				res.add(pr.GreenDop*1000, Command{Tir: tm.Green, Value: 1})
				for _, v := range tm.Reds {
					res.add(pr.GreenDop*1000, Command{Tir: v, Value: 0})
				}

			}
			//Финалочка
			res.add(0, Command{Tir: tm.Yellow, Value: 0})
			res.add(0, Command{Tir: tm.Green, Value: 1})
			for _, v := range tm.Reds {
				res.add(0, Command{Tir: v, Value: 0})
			}
		}
	case 2: //Поворотное направление
		if olds {
			//Нужно закрыть направление Зм Ж Кр
			if pr.GreenBlink > pr.Yellow {
				//Записываем зеленое мигание до желтого
				for i := pr.GreenBlink; i > pr.Yellow; i-- {
					res.add(i*1000, Command{Tir: tm.Green, Value: 0})
					res.add(i*1000-500, Command{Tir: tm.Green, Value: 1})
				}
			}
			if pr.Yellow > pr.Red {
				res.add(pr.Yellow*1000, Command{Tir: tm.Green, Value: 0})
			}
			if pr.Red > 0 {
				res.add(pr.Red*1000, Command{Tir: tm.Green, Value: 0})
			}
			//Финалочка
			res.add(0, Command{Tir: tm.Green, Value: 0})
		} else {
			//Нужно открыть направление Кр КрЖ З
			// if pr.RedYellow != 0 && pr.RedYellow > pr.GreenDop {
			// 	// res.add(pr.RedYellow*1000, Command{Tir: tm.Yellow, Value: 1})
			// 	// for _, v := range tm.Reds {
			// 	// 	res.add(pr.Red*1000, Command{Tir: v, Value: 1})
			// 	// }
			// }
			if pr.GreenDop > 0 {
				res.add(pr.GreenDop*1000, Command{Tir: tm.Green, Value: 1})
			}
			//Финалочка
			res.add(0, Command{Tir: tm.Green, Value: 1})
		}
	case 3: // Пешеходное направление
		if olds {
			//Нужно закрыть направление Зм Ж Кр
			if pr.GreenBlink > pr.Yellow {
				//Записываем зеленое мигание до желтого
				for i := pr.GreenBlink; i > pr.Yellow; i-- {
					res.add(i*1000, Command{Tir: tm.Green, Value: 0})
					res.add(i*1000-500, Command{Tir: tm.Green, Value: 1})
				}
			}
			if pr.Yellow > pr.Red {
				res.add(pr.Yellow*1000, Command{Tir: tm.Green, Value: 0})
			}
			if pr.Red > 0 {
				res.add(pr.Red*1000, Command{Tir: tm.Green, Value: 0})
				for _, v := range tm.Reds {
					res.add(pr.Red*1000, Command{Tir: v, Value: 1})
				}
			}
			//Финалочка
			res.add(0, Command{Tir: tm.Green, Value: 0})
			for _, v := range tm.Reds {
				res.add(0, Command{Tir: v, Value: 1})
			}

		} else {
			//Нужно открыть направление Кр КрЖ З
			// if pr.RedYellow != 0 && pr.RedYellow > pr.GreenDop {
			// 	// res.add(pr.RedYellow*1000, Command{Tir: tm.Yellow, Value: 1})
			// 	// for _, v := range tm.Reds {
			// 	// 	res.add(pr.Red*1000, Command{Tir: v, Value: 1})
			// 	// }
			// }
			if pr.GreenDop > 0 {
				res.add(pr.GreenDop*1000, Command{Tir: tm.Green, Value: 1})
				for _, v := range tm.Reds {
					res.add(pr.GreenDop*1000, Command{Tir: v, Value: 0})
				}

			}
			//Финалочка
			res.add(0, Command{Tir: tm.Green, Value: 1})
			for _, v := range tm.Reds {
				res.add(0, Command{Tir: v, Value: 0})
			}
		}

	}
	return *res
}
