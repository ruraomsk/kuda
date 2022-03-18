package bin

import "fmt"

//GetPromtackt вычисляет переходной промтакт на фазу phase.prom истина если используется базовый промтакт
//	naps текущее состояние направлений Истина если открыто. На выходе изменяем состояние направлений и выдаем
//  план перехода  промтакта который нужно отработать и ошибку
func (c *CMK) GetPromtackt(phase int, prom bool, longtime int) error {
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

	newNaps := make(map[int]bool)
	for n, v := range c.Naps {
		newNaps[n] = v
	}
	for v := range newNaps {
		newNaps[v] = false
	}
	for _, v := range descPhase.Naps {
		newNaps[v] = true
	}
	// fmt.Printf("new %v\nold %v\n", newNaps, c.Naps)
	c.PromMake.new()
	for _, v := range c.TirToNaps {
		fmt.Printf("nap :%d ", v.Number)
		res := v.makeProm(c.Naps[v.Number], newNaps[v.Number], pr[v.Number], longtime)
		for t, k := range res.Ticks {
			for _, l := range k {
				c.PromMake.add(t, l)
			}
		}
		fmt.Println(".")
	}
	for n, v := range newNaps {
		c.Naps[n] = v
	}
	return nil
}
func (tm *TirToNap) makeProm(olds bool, news bool, pr PromTakt, longtime int) PromMake {
	res := new(PromMake)
	res.new()
	if olds == news {
		//Ничего не изменилось
		return *res
	}
	switch tm.Type {
	case 1: //Транспортное направление
		fmt.Print("1 tr ")
		if olds {
			fmt.Print("close ")
			//Нужно закрыть направление Зм Ж Кр
			if pr.GreenBlink > pr.Yellow {
				//Записываем зеленое мигание до желтого
				fmt.Printf("%d GB ", pr.GreenBlink)
				for i := pr.GreenBlink; i > pr.Yellow; i-- {
					res.add(i*1000, Command{Tir: tm.Green, Value: 0})
					res.add(i*1000-500, Command{Tir: tm.Green, Value: 1})
				}
			}
			if pr.Yellow > pr.Red {
				fmt.Printf("%d YE ", pr.Yellow)

				res.add(pr.Yellow*1000, Command{Tir: tm.Green, Value: 0})
				res.add(pr.Yellow*1000, Command{Tir: tm.Yellow, Value: 1})
			}
			if pr.Red >= 0 {
				fmt.Printf("%d RD ", pr.Red)
				res.add(0, Command{Tir: tm.Green, Value: 0})
				res.add(pr.Red*1000, Command{Tir: tm.Yellow, Value: 0})
				for _, v := range tm.Reds {
					res.add(pr.Red*1000, Command{Tir: v, Value: 1})
				}
			}
			if longtime != 0 {

			}
		} else {
			fmt.Print("open ")
			//Нужно открыть направление Кр КрЖ З
			if pr.RedYellow != 0 && pr.RedYellow > pr.GreenDop {
				fmt.Printf("%d RE ", pr.RedYellow)
				res.add(pr.RedYellow*1000, Command{Tir: tm.Yellow, Value: 1})
				for _, v := range tm.Reds {
					res.add(pr.Red*1000, Command{Tir: v, Value: 1})
				}
			}
			if pr.GreenDop > 0 {
				fmt.Printf("%d GD ", pr.GreenDop)

				res.add(pr.GreenDop*1000, Command{Tir: tm.Yellow, Value: 0})
				res.add(pr.GreenDop*1000, Command{Tir: tm.Green, Value: 1})
				for _, v := range tm.Reds {
					res.add(pr.GreenDop*1000, Command{Tir: v, Value: 0})
				}

			} else {
				//Финалочка
				fmt.Printf("%d GR ", 0)
				res.add(0, Command{Tir: tm.Yellow, Value: 0})
				res.add(0, Command{Tir: tm.Green, Value: 1})
				for _, v := range tm.Reds {
					res.add(0, Command{Tir: v, Value: 0})
				}
			}
		}
	case 2: //Поворотное направление
		fmt.Print("2 ro ")
		if olds {
			fmt.Print("close ")
			//Нужно закрыть направление Зм Ж Кр
			if pr.GreenBlink > pr.Yellow {
				fmt.Printf("%d GB ", pr.GreenBlink)
				//Записываем зеленое мигание до желтого
				for i := pr.GreenBlink; i > pr.Yellow; i-- {
					res.add(i*1000, Command{Tir: tm.Green, Value: 0})
					res.add(i*1000-500, Command{Tir: tm.Green, Value: 1})
				}
			}
			if pr.Yellow > pr.Red {
				res.add(pr.Yellow*1000, Command{Tir: tm.Green, Value: 0})
			}
			if pr.Red >= 0 {
				fmt.Printf("%d RD ", pr.Red)
				res.add(pr.Red*1000, Command{Tir: tm.Green, Value: 0})
			}
		} else {
			fmt.Print("open ")
			//Нужно открыть направление Кр КрЖ З
			// if pr.RedYellow != 0 && pr.RedYellow > pr.GreenDop {
			// 	// res.add(pr.RedYellow*1000, Command{Tir: tm.Yellow, Value: 1})
			// 	// for _, v := range tm.Reds {
			// 	// 	res.add(pr.Red*1000, Command{Tir: v, Value: 1})
			// 	// }
			// }
			if pr.GreenDop >= 0 {
				fmt.Printf("%d GB ", pr.GreenDop)
				res.add(pr.GreenDop*1000, Command{Tir: tm.Green, Value: 1})
			}
		}
	case 3: // Пешеходное направление
		fmt.Print("3 st ")
		if olds {
			fmt.Print("close ")
			//Нужно закрыть направление Зм Ж Кр
			if pr.GreenBlink > pr.Yellow {
				//Записываем зеленое мигание до желтого
				fmt.Printf("%d GB ", pr.GreenBlink)
				for i := pr.GreenBlink; i > pr.Yellow; i-- {
					res.add(i*1000, Command{Tir: tm.Green, Value: 0})
					res.add(i*1000-500, Command{Tir: tm.Green, Value: 1})
				}
			}
			if pr.Yellow > pr.Red {
				fmt.Printf("%d YE ", pr.Yellow)
				res.add(pr.Yellow*1000, Command{Tir: tm.Green, Value: 0})
			}
			if pr.Red >= 0 {
				fmt.Printf("%d RD ", pr.Red)
				res.add(pr.Red*1000, Command{Tir: tm.Green, Value: 0})
				for _, v := range tm.Reds {
					res.add(pr.Red*1000, Command{Tir: v, Value: 1})
				}
			}

		} else {
			fmt.Print("open ")
			//Нужно открыть направление Кр КрЖ З
			// if pr.RedYellow != 0 && pr.RedYellow > pr.GreenDop {
			// 	// res.add(pr.RedYellow*1000, Command{Tir: tm.Yellow, Value: 1})
			// 	// for _, v := range tm.Reds {
			// 	// 	res.add(pr.Red*1000, Command{Tir: v, Value: 1})
			// 	// }
			// }
			if pr.GreenDop >= 0 {
				fmt.Printf("%d GD ", pr.GreenDop)
				res.add(pr.GreenDop*1000, Command{Tir: tm.Green, Value: 1})
				for _, v := range tm.Reds {
					res.add(pr.GreenDop*1000, Command{Tir: v, Value: 0})
				}
			}
		}

	}
	// fmt.Printf(" %v", res)
	return *res
}
