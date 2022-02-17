package hardware

var (
	Cpu = ModuleCPU{
		moduleNumber:  1,
		moduleSlaveID: 255,
		moduleType:    0,
		moduleStatus:  1,
		moduleSetup:   38,
		size:          549,
		di:            Dis,
		do:            Dos,
		ai:            Ais,
		ao:            Aos,
	}

	Dis = map[int]DI{
		1:  {value: bs{w: 2, b: 0}, counter: 4, state: 39, reset: bs{w: 75, b: 0}, bounce: 500, blocked: 520, bl_state: bs{w: 540, b: 0}, front: bs{w: 542, b: 0}},
		2:  {value: bs{w: 2, b: 1}, counter: 5, state: 40, reset: bs{w: 75, b: 1}, bounce: 501, blocked: 521, bl_state: bs{w: 540, b: 1}, front: bs{w: 542, b: 1}},
		3:  {value: bs{w: 2, b: 2}, state: 41, reset: bs{w: 75, b: 2}, bounce: 502, blocked: 522, bl_state: bs{w: 540, b: 2}, front: bs{w: 542, b: 2}},
		4:  {value: bs{w: 2, b: 3}, state: 42, reset: bs{w: 75, b: 3}, bounce: 503, blocked: 523, bl_state: bs{w: 540, b: 3}, front: bs{w: 542, b: 3}},
		5:  {value: bs{w: 2, b: 4}, state: 43, reset: bs{w: 75, b: 4}, bounce: 504, blocked: 524, bl_state: bs{w: 540, b: 4}, front: bs{w: 542, b: 4}},
		6:  {value: bs{w: 2, b: 5}, state: 44, reset: bs{w: 75, b: 5}, bounce: 505, blocked: 525, bl_state: bs{w: 540, b: 5}, front: bs{w: 542, b: 5}},
		7:  {value: bs{w: 2, b: 6}, state: 45, reset: bs{w: 75, b: 6}, bounce: 506, blocked: 526, bl_state: bs{w: 540, b: 6}, front: bs{w: 542, b: 6}},
		8:  {value: bs{w: 2, b: 7}, state: 46, reset: bs{w: 75, b: 7}, bounce: 507, blocked: 527, bl_state: bs{w: 540, b: 7}, front: bs{w: 542, b: 7}},
		9:  {value: bs{w: 2, b: 8}, state: 47, reset: bs{w: 75, b: 8}, bounce: 508, blocked: 528, bl_state: bs{w: 540, b: 8}, front: bs{w: 542, b: 8}},
		10: {value: bs{w: 2, b: 9}, state: 48, reset: bs{w: 75, b: 9}, bounce: 509, blocked: 529, bl_state: bs{w: 540, b: 9}, front: bs{w: 542, b: 9}},
		11: {value: bs{w: 2, b: 10}, state: 49, reset: bs{w: 75, b: 10}, bounce: 510, blocked: 530, bl_state: bs{w: 540, b: 10}, front: bs{w: 542, b: 10}},
		12: {value: bs{w: 2, b: 11}, state: 50, reset: bs{w: 75, b: 11}, bounce: 511, blocked: 531, bl_state: bs{w: 540, b: 11}, front: bs{w: 542, b: 11}},
		13: {value: bs{w: 2, b: 12}, state: 51, reset: bs{w: 75, b: 12}, bounce: 512, blocked: 532, bl_state: bs{w: 540, b: 12}, front: bs{w: 542, b: 12}},
		14: {value: bs{w: 2, b: 13}, state: 52, reset: bs{w: 75, b: 13}, bounce: 513, blocked: 533, bl_state: bs{w: 540, b: 13}, front: bs{w: 542, b: 13}},
		15: {value: bs{w: 2, b: 14}, state: 53, reset: bs{w: 75, b: 14}, bounce: 514, blocked: 534, bl_state: bs{w: 540, b: 14}, front: bs{w: 542, b: 14}},
		16: {value: bs{w: 2, b: 15}, state: 54, reset: bs{w: 75, b: 15}, bounce: 515, blocked: 535, bl_state: bs{w: 540, b: 15}, front: bs{w: 542, b: 15}},
		17: {value: bs{w: 3, b: 0}, state: 55, reset: bs{w: 76, b: 0}, bounce: 516, blocked: 536, bl_state: bs{w: 541, b: 0}, front: bs{w: 543, b: 0}},
		18: {value: bs{w: 3, b: 1}, state: 56, reset: bs{w: 76, b: 1}, bounce: 517, blocked: 537, bl_state: bs{w: 541, b: 1}, front: bs{w: 543, b: 1}},
		19: {value: bs{w: 3, b: 2}, state: 57, reset: bs{w: 76, b: 2}, bounce: 518, blocked: 538, bl_state: bs{w: 541, b: 2}, front: bs{w: 543, b: 2}},
		20: {value: bs{w: 3, b: 3}, state: 58, reset: bs{w: 76, b: 3}, bounce: 519, blocked: 539, bl_state: bs{w: 541, b: 3}, front: bs{w: 543, b: 3}},
	}
	Dos = map[int]DO{
		1:  {value: bs{w: 78, b: 0}, state: 63},
		2:  {value: bs{w: 78, b: 1}, state: 64},
		3:  {value: bs{w: 78, b: 2}, state: 65, kz: bs{w: 36, b: 2}, kz_sbros: bs{w: 77, b: 2}, kz_ctrl: bs{w: 544, b: 2}},
		4:  {value: bs{w: 78, b: 3}, state: 66, kz: bs{w: 36, b: 3}, kz_sbros: bs{w: 77, b: 3}, kz_ctrl: bs{w: 544, b: 3}},
		5:  {value: bs{w: 78, b: 4}, state: 67, kz: bs{w: 36, b: 4}, kz_sbros: bs{w: 77, b: 4}, kz_ctrl: bs{w: 544, b: 4}},
		6:  {value: bs{w: 78, b: 5}, state: 68, kz: bs{w: 36, b: 5}, kz_sbros: bs{w: 77, b: 5}, kz_ctrl: bs{w: 544, b: 5}},
		7:  {value: bs{w: 78, b: 6}, state: 69, kz: bs{w: 36, b: 6}, kz_sbros: bs{w: 77, b: 6}, kz_ctrl: bs{w: 544, b: 6}},
		8:  {value: bs{w: 78, b: 7}, state: 70, kz: bs{w: 36, b: 7}, kz_sbros: bs{w: 77, b: 7}, kz_ctrl: bs{w: 544, b: 7}},
		9:  {value: bs{w: 78, b: 8}, state: 71, kz: bs{w: 36, b: 8}, kz_sbros: bs{w: 77, b: 8}, kz_ctrl: bs{w: 544, b: 8}},
		10: {value: bs{w: 78, b: 9}, state: 72, kz: bs{w: 36, b: 9}, kz_sbros: bs{w: 77, b: 9}, kz_ctrl: bs{w: 544, b: 9}},
	}
	Ais = map[int]AI{
		1: {value: 24, fvalue: 28, filter: 546},
		2: {value: 25, fvalue: 30, filter: 547},
		3: {value: 26, fvalue: 32, filter: 548},
		4: {value: 27, fvalue: 34, filter: 549},
	}
	Aos = map[int]AO{
		1: {value: 98, fvalue: 100, state: 73},
		2: {value: 99, fvalue: 102, state: 74},
	}
)
