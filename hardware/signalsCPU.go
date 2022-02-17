package hardware

import "sync"

type ModuleCPU struct {
	moduleSlaveID int
	moduleNumber  int
	moduleType    int //Место хранения номера модуля
	moduleStatus  int //Место хранения статуса модуля
	moduleSetup   int //Место хранения настройки модуля
	size          int //Размер буфера чтения
	di            map[int]DI
	ai            map[int]AI
	do            map[int]DO
	ao            map[int]AO
	masterTCP     *MasterTcp
	work          bool
	mutex         sync.Mutex
	writer        chan writeHR
}

type bs struct {
	w int //номер слова
	b int //номер бита
}
type writeHR struct {
	pos   bs
	value int
}
type DI struct {
	value    bs  //Место хранения значения
	counter  int //Место хранения счетчика
	state    int //Место хранения/установки режима работы
	reset    bs  //Место сброса состояний в режимах Вход с удержанием Тригер и Счетчик
	bounce   int //Место хранения времени антидребезга
	blocked  int //Место хранения времени удержания входа
	bl_state bs  //Место удерживаемых состояний
	front    bs  //Место хранения активных фронтов
}

func (d *DI) getValue(hr []uint16) bool {
	v := hr[d.value.w]
	v = v >> uint16(d.value.b)
	return v&1 == 1
}

type AI struct {
	value  int //Место хранения код АЦП значения
	fvalue int //Место хранения float значения
	filter int //Место хранения фильтра низких частот
}

type DO struct {
	value    bs  //Место хранения значения
	state    int //Место хранения режима работы
	kz       bs  //Место хранения состояния КЗ
	kz_sbros bs  //Место сброса КЗ выхода
	kz_ctrl  bs  //Место установки защиты от КЗ
}

func (d *DO) getValue(hr []uint16) bool {
	v := hr[d.value.w]
	v = v >> uint16(d.value.b)
	return v&1 == 1
}

type AO struct {
	value  int //Место хранения код АЦП значения
	fvalue int //Место хранения float значения
	state  int //Место хранения режима работы
}
