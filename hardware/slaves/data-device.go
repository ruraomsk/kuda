package slaves

import (
	"fmt"
	"time"

	"github.com/tbrandon/mbserver"
)

/*
	0 	Тип модуля 			чтение	номер модуля
	1	Статус модуля 		чтение	бит 0 - работа
	2	Настройка модуля    чтение
							запись
	3   Состояние выходов 	чтение	Бит 0 канал 1 /значение 1 включен 0 выключен
							запись	Бит 7 канал 8
	4	Состояние КЗ 		чтение  Бит 0 канал 1 /значение 1 есть КЗ 0 норма
	 								Бит 7 канал 8
	5	Состояние Обрыв		чтение	Бит 0 канал 1 /значение 1 есть Обрыв 0 норма
	 								Бит 7 канал 8
	6	Несанкц включение	чтение 	Бит 0 канал 1 /значение 1 есть Вкл 0 норма
	 								Бит 7 канал 8

*/
const (
	value  = 3
	kz     = 4
	broken = 5
	bad    = 6
)

type DeviceC8 struct {
	Port   int
	Number int
	Size   int
	server *mbserver.Server
}

func (d *DeviceC8) WorkDevice() {
	d.server = mbserver.NewServer()
	con := fmt.Sprintf(":%d", d.Port)
	for i := 0; i < d.Size; i++ {
		d.server.HoldingRegisters[i] = 0
	}
	d.server.HoldingRegisters[0] = uint16(d.Number)
	d.server.HoldingRegisters[1] = 1
	_ = d.server.ListenTCP(con)
	fmt.Printf("device %d ready\n", d.Number)
	for {
		time.Sleep(time.Second)
	}
}
func (d *DeviceC8) GetValue(can int) bool {
	if can < 1 || can > 8 {
		return false
	}
	mask := 1
	for i := 1; i < can; i++ {
		mask = mask << 1
	}
	return (d.server.HoldingRegisters[value] & uint16(mask)) == 1
}
func (d *DeviceC8) GetKZ(can int) bool {
	if can < 1 || can > 8 {
		return false
	}
	mask := 1
	for i := 1; i < can; i++ {
		mask = mask << 1
	}
	return (d.server.HoldingRegisters[kz] & uint16(mask)) == 1
}
func (d *DeviceC8) GetBrocken(can int) bool {
	if can < 1 || can > 8 {
		return false
	}
	mask := 1
	for i := 1; i < can; i++ {
		mask = mask << 1
	}
	return (d.server.HoldingRegisters[broken] & uint16(mask)) == 1
}
func (d *DeviceC8) GetBad(can int) bool {
	if can < 1 || can > 8 {
		return false
	}
	mask := 1
	for i := 1; i < can; i++ {
		mask = mask << 1
	}
	return (d.server.HoldingRegisters[bad] & uint16(mask)) == 1
}
