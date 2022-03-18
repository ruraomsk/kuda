package bin

//GetBaseOrUniver возвращает тип рекомендованного промтакта
func (c *CMK) GetBaseOrUniver(oldPhase, newPhase int) bool {
	//Находим место old
	for o, v := range c.RPUs[0].Phases {
		if v.Phase == oldPhase {
			//Смотрим слева
			o--
			if o < 0 {
				o = len(c.RPUs[0].Phases) - 1
			}
			if c.RPUs[0].Phases[o].Phase == newPhase {
				return true
			}
			//Смотрим справа
			o++
			if o >= len(c.RPUs[0].Phases) {
				o = 0
			}
			if c.RPUs[0].Phases[o].Phase == newPhase {
				return true
			}

		}
	}
	return false
}
func (c *CMK) IsPhase(phase int) bool {
	for _, v := range c.NtoPhases {
		if v.NumPhase == phase {
			return true
		}
	}
	return false
}
