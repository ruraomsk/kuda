package hardware

import "sync"

type ModuleC8 struct {
	moduleSlaveID int
	moduleNumber  int
	moduleType    int //Место хранения номера модуля
	moduleStatus  int //Место хранения статуса модуля
	moduleSetup   int //Место хранения настройки модуля
	size          int //Размер буфера чтения
	connect       string
	с8            map[int]C8
	masterTCP     *MasterTcp
	work          bool
	mutex         sync.Mutex
	writer        chan writeHR
}
type C8 struct {
	value   bs //Место хранения значения
	kz      bs //Место хранения состояния КЗ
	brocken bs //Место хранения состояния обрыва
	bad     bs //Место хранения признака несанционированного включения
}

func (d *C8) getValue(hr []uint16) bool {
	v := hr[d.value.w]
	v = v >> uint16(d.value.b)
	return v&1 == 1
}
