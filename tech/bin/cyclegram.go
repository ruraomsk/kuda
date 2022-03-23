package bin

import (
	"fmt"
)

//
//makeCycleGramm Строим циклограмму всего плана РПУ plan номер плана
func (c *CMK) MakeCycleGramm(plan int) error {
	saveNaps := make(map[int]bool)
	for n, v := range c.Naps {
		saveNaps[n] = v
	}
	lensProm := make(map[int]int)
	step := 0
	for step < len(c.RPUs[0].Phases) {
		err := cmk.GetPromtackt(c.RPUs[0].Phases[step].Phase, true)
		if err != nil {
			return err
		}

		timeCount = cmk.PromMake.GetMaxTime()
		lensProm[c.RPUs[0].Phases[step].Phase] = timeCount / 1000
		step++
	}
	// fmt.Printf("lens %v\n", lensProm)
	for n, v := range saveNaps {
		c.Naps[n] = v
	}

	for _, nap := range c.TirToNaps {
		if nap.Counter == 0 {
			continue
		}

		allClose := true
		oldNaps := make(map[int]bool)
		for n, v := range c.Naps {
			oldNaps[n] = v
			if v {
				allClose = false
			}
		}
		fmt.Printf("naps %v %v\n", c.Naps, allClose)
		timeNow := TimeNowOfSecond()
		//Удалим все счетчики старше начала
		a := make([]int, 0)
		for k := range counts {
			if k > timeNow {
				a = append(a, k)
			}
		}
		for _, v := range a {
			delete(counts, v)
		}
		startTime := -1
		if allClose {
			startTime = timeNow
		}
		for i := 0; i < 3; i++ {
			step = 0
			for step < len(c.RPUs[0].Phases) {
				phase := c.RPUs[0].Phases[step].Phase
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
				prom := c.PrMapBase[nap.Number]
				lpr := lensProm[phase]
				newNaps := makeNewNaps(oldNaps, descPhase)
				fmt.Printf("timenow %d %d %d %d %v", timeNow, c.RPUs[0].Phases[step].Time, phase, lpr, newNaps)
				if newNaps[nap.Number] == oldNaps[nap.Number] {
					timeNow += c.RPUs[0].Phases[step].Time
					oldNaps = moveNaps(newNaps)
					fmt.Printf("dont change \n")
					step++

					continue
				}
				switch nap.Type {
				case 1: //Транспортное
					if oldNaps[nap.Number] {
						//Закрываем направление в начале запишем время когда ушел зеленый
						fmt.Printf("tr close ")
						green := 0
						if prom.Yellow > prom.Red {
							green = prom.Yellow
						} else {
							if prom.Red >= 0 {
								green = prom.Red

							}
						}
						if green > lpr {
							green = 0
						}
						if startTime > 0 {
							fmt.Printf("green count start %d lenght %d", startTime, (timeNow+lpr-green)-startTime)
							setCounter(startTime, nap.Counter, (timeNow+lpr-green)-startTime)
						}
						startTime = timeNow + lpr - prom.Red
						if prom.Red > lpr {
							startTime = timeNow + lpr
						}
						startTime++
					} else {
						//Открываем направление в начале запишем время когда ушел красный
						fmt.Printf("tr open ")
						red := 0
						red = prom.GreenDop
						if red > lpr {
							red = 0
						}
						if startTime > 0 {
							fmt.Printf("red count start %d lenght %d", startTime, (timeNow+lpr-red)-startTime)
							setCounter(startTime, nap.Counter, (timeNow+lpr-red)-startTime)
						}
						startTime = timeNow + lpr - red
						startTime++
					}
				case 2: //Поворотное
					if oldNaps[nap.Number] {
						//Закрываем направление
						fmt.Printf("pr close ")

					} else {
						//Открываем направление
						fmt.Printf("pr open ")
					}
				case 3: //Пешеходное
					if oldNaps[nap.Number] {
						//Закрываем направление
						fmt.Printf("st close ")
					} else {
						fmt.Printf("st open ")
						//Открываем направление
					}
				}
				fmt.Println(".")
				timeNow += c.RPUs[0].Phases[step].Time
				oldNaps = moveNaps(newNaps)
				step++
			}

		}

	}

	return nil
}

func makeNewNaps(oldNaps map[int]bool, phase NtoPhase) (newNaps map[int]bool) {
	newNaps = make(map[int]bool)
	for n := range oldNaps {
		newNaps[n] = false
	}
	for _, v := range phase.Naps {
		newNaps[v] = true
	}

	return
}
func moveNaps(newNaps map[int]bool) (oldNaps map[int]bool) {
	oldNaps = make(map[int]bool)
	for n, v := range newNaps {
		oldNaps[n] = v
	}
	return
}
