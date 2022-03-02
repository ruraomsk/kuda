package bin

import (
	"sort"

	"github.com/ruraomsk/ag-server/logger"
)

func (c *CMK) convertToStruct(ptr, lenght, count int) {
	logger.Debug.Printf("getMass %d %x %d %d", nm, ptr, lenght, count)
	switch nm {
	case 1:
		//
		c.NtoPhases = make([]NtoPhase, 0)
		for i := 0; i < count; i++ {
			n := NtoPhase{NumPhase: i + 1, Naps: make([]int, 0)}
			for j := 0; j < lenght; j++ {
				if buff[ptr] != 255 {
					n.Naps = append(n.Naps, int(buff[ptr]))
				}
				ptr++
			}
			c.NtoPhases = append(c.NtoPhases, n)
		}
		sort.Slice(c.NtoPhases, func(i, j int) bool {
			return c.NtoPhases[i].NumPhase < c.NtoPhases[j].NumPhase
		})
		return
	case 2:
		c.TirToNaps = make([]TirToNap, 0)
		for i := 0; i < count; i++ {
			t := TirToNap{Number: i + 1, Green: int(buff[ptr]), Yellow: int(buff[ptr+1]), Reds: make([]int, 0)}
			ptr += 2
			for j := 0; j < 3; j++ {
				if buff[ptr] != 255 {
					t.Reds = append(t.Reds, int(buff[ptr]))
				}
				ptr++
			}
			c.TirToNaps = append(c.TirToNaps, t)
		}
		sort.Slice(c.TirToNaps, func(i, j int) bool {
			return c.TirToNaps[i].Number < c.TirToNaps[j].Number
		})
		return
	case 3:
		c.TminToPhases = make([]TimeToPhase, 0)
		for i := 0; i < count; i++ {
			if buff[ptr] != 255 {
				c.TminToPhases = append(c.TminToPhases, TimeToPhase{NumPhase: int(buff[ptr]), Tmin: int(buff[ptr+1])})
			}
			ptr += 2
		}
		return
	case 4:
		c.RPUs = make([]RPU, 0)
		for i := 0; i < count; i++ {
			r := RPU{Number: i + 1, Tcycle: int(buff[ptr]), Continue: buff[ptr+3] == 1, Phases: make([]Phase, 0)}
			ptr += 4
			for j := 4; j < lenght; j += 4 {

				if buff[ptr] != 255 {
					r.Phases = append(r.Phases, Phase{TVP1: int(buff[ptr]), TVP2: int(buff[ptr+1]), ZamPh: int(buff[ptr+2]), Time: int(buff[ptr+3])})
				}
				ptr += 4
			}
			c.RPUs = append(c.RPUs, r)
		}
		sort.Slice(c.RPUs, func(i, j int) bool {
			return c.RPUs[i].Number < c.RPUs[j].Number
		})
		return
	case 5:
		c.PromTaktBases = make([]PromTakt, 0)
		for i := 0; i < count; i++ {
			p := PromTakt{Nap: i + 1, GreenDop: int(buff[ptr]), GreenBlink: int(buff[ptr+1]), Yellow: int(buff[ptr+2]), Red: int(buff[ptr+3]), RedYellow: int(buff[ptr+4])}
			c.PromTaktBases = append(c.PromTaktBases, p)
			ptr += 5
		}
		sort.Slice(c.PromTaktBases, func(i, j int) bool {
			return c.PromTaktBases[i].Nap < c.PromTaktBases[j].Nap
		})
		return
	case 6:
		c.PromTakt = make([]PromTakt, 0)
		for i := 0; i < count; i++ {
			p := PromTakt{Nap: i + 1, GreenDop: int(buff[ptr]), GreenBlink: int(buff[ptr+1]), Yellow: int(buff[ptr+2]), Red: int(buff[ptr+3]), RedYellow: int(buff[ptr+4])}
			c.PromTakt = append(c.PromTakt, p)
			ptr += 5
		}
		sort.Slice(c.PromTakt, func(i, j int) bool {
			return c.PromTakt[i].Nap < c.PromTakt[j].Nap
		})
		return
	case 7:
		c.RedGroups = make([]Mask, 0)
		for i := 0; i < count; i++ {
			m := Mask{Number: i + 1, Mask: make([]byte, 0)}
			for j := 0; j < lenght; j++ {
				m.Mask = append(m.Mask, buff[ptr])
				ptr++
			}
			c.RedGroups = append(c.RedGroups, m)
		}
		sort.Slice(c.RedGroups, func(i, j int) bool {
			return c.RedGroups[i].Number < c.RedGroups[j].Number
		})
		return
	case 8:
		for i := 0; i < len(c.NKs); i++ {
			c.NKs[i] = int(buff[ptr])
			ptr++
		}
		return
	case 9:
		c.CKs = make([]CK, 0)
		for i := 0; i < count; i++ {
			ck := CK{Number: i + 1, Lines: make([]Line, 0)}
			for j := 0; j < lenght; j++ {
				if buff[ptr] != 255 {
					ck.Lines = append(ck.Lines, Line{Number: int(buff[ptr]), Time: (int(buff[ptr+1]) * 60) + int(buff[ptr+2])})
				}
				ptr += 3
			}
			c.CKs = append(c.CKs, ck)
		}
		sort.Slice(c.CKs, func(i, j int) bool {
			return c.CKs[i].Number < c.CKs[j].Number
		})
		return
	case 10:
		c.Konf = make([]Mask, 0)
		for i := 0; i < count; i++ {
			m := Mask{Number: i + 1, Mask: make([]byte, 0)}
			for j := 0; j < lenght; j++ {
				m.Mask = append(m.Mask, buff[ptr])
				ptr++
			}
			c.Konf = append(c.Konf, m)
		}
		sort.Slice(c.Konf, func(i, j int) bool {
			return c.Konf[i].Number < c.Konf[j].Number
		})
		return
	case 11:
		c.TVPs = make([]TVP, 0)
		for i := 0; i < count; i++ {
			m := 0
			for j := 0; j < 8; j++ {
				if buff[ptr+1]>>j&1 == 1 {
					m = j
					break
				}
			}
			tvp := TVP{Number: i + 1, Wait: int(buff[ptr])*8 + m, Phases: make([]int, 0)}
			ptr += 2
			for j := 2; j < lenght; j++ {
				if buff[ptr] != 255 {
					tvp.Phases = append(tvp.Phases, int(buff[ptr]))
				}
				ptr++
			}
			c.TVPs = append(c.TVPs, tvp)
		}
		sort.Slice(c.TVPs, func(i, j int) bool {
			return c.TVPs[i].Number < c.TVPs[j].Number
		})
		return
	case 12:
		c.GreenBlinks = make([]NtoPhase, 0)
		for i := 0; i < count; i++ {
			n := NtoPhase{NumPhase: i + 1, Naps: make([]int, 0)}
			for j := 0; j < lenght; j++ {
				if buff[ptr] != 255 {
					n.Naps = append(n.Naps, int(buff[ptr]))
				}
				ptr++
			}
			c.GreenBlinks = append(c.GreenBlinks, n)
		}
		sort.Slice(c.GreenBlinks, func(i, j int) bool {
			return c.GreenBlinks[i].NumPhase < c.GreenBlinks[j].NumPhase
		})
		return
	case 13:
		c.TmaxToPhases = make([]TimeToPhase, 0)
		for i := 0; i < count*lenght; i++ {
			if buff[ptr] == 255 {
				break
			}
			c.TmaxToPhases = append(c.TmaxToPhases, TimeToPhase{NumPhase: int(buff[ptr]), Tmin: int(buff[ptr+1])})
			ptr += 2
		}
		return
	case 14:
		c.GreenToOut = make([]Mask, 0)
		for i := 0; i < count; i++ {
			m := Mask{Number: i + 1, Mask: make([]byte, 0)}
			for j := 0; j < lenght; j++ {
				m.Mask = append(m.Mask, buff[ptr])
				ptr++
			}
			c.GreenToOut = append(c.GreenToOut, m)
		}
		sort.Slice(c.GreenToOut, func(i, j int) bool {
			return c.GreenToOut[i].Number < c.GreenToOut[j].Number
		})
		return
	case 15:
		c.PowerToOut = make([]Mask, 0)
		for i := 0; i < count; i++ {
			m := Mask{Number: i + 1, Mask: make([]byte, 0)}
			for j := 0; j < lenght; j++ {
				m.Mask = append(m.Mask, buff[ptr])
				ptr++
			}
			c.PowerToOut = append(c.PowerToOut, m)
		}
		sort.Slice(c.PowerToOut, func(i, j int) bool {
			return c.PowerToOut[i].Number < c.PowerToOut[j].Number
		})
		return
	case 16:
		c.ReedToOut = make([]Mask, 0)
		for i := 0; i < count; i++ {
			m := Mask{Number: i + 1, Mask: make([]byte, 0)}
			for j := 0; j < lenght; j++ {
				m.Mask = append(m.Mask, buff[ptr])
				ptr++
			}
			c.ReedToOut = append(c.ReedToOut, m)
		}
		sort.Slice(c.ReedToOut, func(i, j int) bool {
			return c.ReedToOut[i].Number < c.ReedToOut[j].Number
		})
		return
	case 17:
		c.YellowToOut = make([]Mask, 0)
		for i := 0; i < count; i++ {
			m := Mask{Number: i + 1, Mask: make([]byte, 0)}
			for j := 0; j < lenght; j++ {
				m.Mask = append(m.Mask, buff[ptr])
				ptr++
			}
			c.YellowToOut = append(c.YellowToOut, m)
		}
		sort.Slice(c.YellowToOut, func(i, j int) bool {
			return c.YellowToOut[i].Number < c.YellowToOut[j].Number
		})
		return
	}

}
