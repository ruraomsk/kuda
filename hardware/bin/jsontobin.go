package bin

func (c *CMK) toBin(n int) (buffer []byte, lenght, count int) {
	buffer = make([]byte, 0)
	lenght = 0
	count = 0
	switch n {
	case 1:
		count = len(c.NtoPhases)
		for _, v := range c.NtoPhases {
			if len(v.Naps) > lenght {
				lenght = len(v.Naps)
			}
		}
		lenght++
		for _, v := range c.NtoPhases {
			for i := 0; i < lenght; i++ {
				if i < len(v.Naps) {
					buffer = append(buffer, byte(v.Naps[i]))
				} else {
					buffer = append(buffer, 255)
				}
			}
		}
		return
	case 2:
		count = len(c.TirToNaps)
		for _, v := range c.TirToNaps {
			if len(v.Reds) > lenght {
				lenght = len(v.Reds)
			}
		}
		for _, v := range c.TirToNaps {
			buffer = append(buffer, byte(v.Green))
			buffer = append(buffer, byte(v.Yellow))
			for i := 0; i < lenght; i++ {
				if i < len(v.Reds) {
					buffer = append(buffer, byte(v.Reds[i]))
				} else {
					buffer = append(buffer, 255)
				}
			}
			buffer = append(buffer, 255)
		}
		lenght += 3
		return
	case 3:
		count = 1
		lenght = len(c.TminToPhases) + 1
		for _, v := range c.TminToPhases {
			buffer = append(buffer, byte(v.NumPhase))
			buffer = append(buffer, byte(v.Tmin))
		}
		buffer = append(buffer, 255)
		return
	case 4:
		count = len(c.RPUs)
		for _, v := range c.RPUs {
			if len(v.Phases) > lenght {
				lenght = len(v.Phases)
			}
		}
		for _, v := range c.RPUs {
			buffer = append(buffer, byte(v.Tcycle))
			buffer = append(buffer, 0)
			buffer = append(buffer, 0)
			if v.Continue {
				buffer = append(buffer, 1)
			} else {
				buffer = append(buffer, 0)
			}
			for i := 0; i < lenght; i++ {
				if i < len(v.Phases) {
					buffer = append(buffer, byte(v.Phases[i].TVP1))
					buffer = append(buffer, byte(v.Phases[i].TVP2))
					buffer = append(buffer, byte(v.Phases[i].ZamPh))
					buffer = append(buffer, byte(v.Phases[i].Time))
				} else {
					buffer = append(buffer, 255)
					buffer = append(buffer, 255)
					buffer = append(buffer, 255)
					buffer = append(buffer, 255)
				}
			}
			buffer = append(buffer, 255)
		}
		lenght = 4 + 4*lenght + 1
		return
	case 5:
		count = len(c.PromTaktBases)
		lenght = 5
		for _, v := range c.PromTaktBases {
			buffer = append(buffer, byte(v.GreenDop))
			buffer = append(buffer, byte(v.GreenBlink))
			buffer = append(buffer, byte(v.Yellow))
			buffer = append(buffer, byte(v.Red))
			buffer = append(buffer, byte(v.RedYellow))
		}
		return
	case 6:
		count = len(c.PromTakt)
		lenght = 5
		for _, v := range c.PromTakt {
			buffer = append(buffer, byte(v.GreenDop))
			buffer = append(buffer, byte(v.GreenBlink))
			buffer = append(buffer, byte(v.Yellow))
			buffer = append(buffer, byte(v.Red))
			buffer = append(buffer, byte(v.RedYellow))
		}
		return
	case 7:
		count = len(c.RedGroups)
		lenght = len(c.RedGroups[0].Mask)
		for _, v := range c.RedGroups {
			buffer = append(buffer, v.Mask...)
		}
		return
	case 8:
		count = 1
		lenght = 7
		for _, v := range c.NKs {
			buffer = append(buffer, byte(v))
		}
		return
	case 9:
		count = len(c.CKs)
		for _, v := range c.CKs {
			if len(v.Lines) > lenght {
				lenght = len(v.Lines)
			}
		}
		for _, v := range c.CKs {
			for i := 0; i < lenght; i++ {
				if i < len(v.Lines) {
					buffer = append(buffer, byte(v.Lines[i].Number))
					buffer = append(buffer, byte(v.Lines[i].Time/60))
					buffer = append(buffer, byte(v.Lines[i].Time%60))
				} else {
					buffer = append(buffer, 255)
					buffer = append(buffer, 255)
					buffer = append(buffer, 255)
				}
			}
			buffer = append(buffer, 255)
		}
		lenght = lenght*3 + 1
		return
	case 10:
		count = len(c.Konf)
		lenght = len(c.Konf[0].Mask)
		for _, v := range c.Konf {
			buffer = append(buffer, v.Mask...)
		}
		return
	case 11:
		count = len(c.TVPs)
		for _, v := range c.TVPs {
			if len(v.Phases) > lenght {
				lenght = len(v.Phases)
			}
		}
		for _, v := range c.TVPs {
			buffer = append(buffer, byte(v.Wait/8))
			n := 1
			for i := 0; i < v.Wait%8; i++ {
				n = n << 1
			}
			buffer = append(buffer, byte(n))
			for i := 0; i < lenght; i++ {
				if i < len(v.Phases) {
					buffer = append(buffer, byte(v.Phases[i]))
				} else {
					buffer = append(buffer, 255)
				}
			}
			buffer = append(buffer, 255)
		}
		lenght = lenght + 3
		return
	case 12:
		count = len(c.GreenBlinks)
		for _, v := range c.GreenBlinks {
			if len(v.Naps) > lenght {
				lenght = len(v.Naps)
			}
		}
		for _, v := range c.GreenBlinks {
			for i := 0; i < lenght; i++ {
				if i < len(v.Naps) {
					buffer = append(buffer, byte(v.Naps[i]))
				} else {
					buffer = append(buffer, 255)
				}
			}
			buffer = append(buffer, 255)
		}
		lenght++
		return
	case 13:
		count = 1
		for _, v := range c.TmaxToPhases {
			buffer = append(buffer, byte(v.NumPhase))
			buffer = append(buffer, byte(v.Tmin))
		}
		buffer = append(buffer, 255)
		lenght = len(c.TmaxToPhases)*2 + 1
		return
	case 14:
		count = 1
		lenght = len(c.GreenToOut)
		for _, v := range c.GreenToOut {
			buffer = append(buffer, v.Mask...)
		}
		return
	case 15:
		count = 1
		lenght = len(c.PowerToOut)
		for _, v := range c.PowerToOut {
			buffer = append(buffer, v.Mask...)
		}
		return
	case 16:
		count = 1
		lenght = len(c.ReedToOut)
		for _, v := range c.ReedToOut {
			buffer = append(buffer, v.Mask...)
		}
		return
	case 17:
		count = 1
		lenght = len(c.YellowToOut)
		for _, v := range c.YellowToOut {
			buffer = append(buffer, v.Mask...)
		}
		return
		// case 18:
		// 	count = 1
		// 	lenght = len(c.Phases)
		// 	for _, v := range c.Phases {
		// 		buffer = append(buffer, byte(v))
		// 	}
		// 	buffer = append(buffer, 255)
		// 	lenght++
		// 	return
		// case 19:
		// 	count = 1
		// 	lenght = len(c.Switches)
		// 	for _, v := range c.Switches {
		// 		buffer = append(buffer, byte(v))
		// 	}
		// 	buffer = append(buffer, 255)
		// 	lenght++
		// 	return

	}
	buffer = append(buffer, 255)
	lenght = 1
	count = 1
	return
}
