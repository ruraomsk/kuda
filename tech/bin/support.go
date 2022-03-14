package bin

//GetBaseOrUniver возвращает тип рекомендованного промтакта
func (c *CMK) GetBaseOrUniver(oldPhase, newPhase int) bool {
	//Находим место old
	for o, v := range c.RPUs {
		if v.Number == oldPhase {
			//Смотрим слева
			o--
			if o < 0 {
				o = len(c.RPUs)
			}
			if c.RPUs[o].Number == newPhase {
				return true
			}
			//Смотрим справа
			o++
			if o == len(c.RPUs) {
				o = 0
			}
			if c.RPUs[o].Number == newPhase {
				return true
			}

		}
	}
	return false
}
