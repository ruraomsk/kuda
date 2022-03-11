package bin

type CMK struct {
	BaseOption    int              `json:"bopt"`         // Базовые опции
	BaseRPU       int              `json:"brpu"`         // Номер базовой РПУ
	NumTU         int              `json:"ntu"`          // Номер ТУ
	NumDK         int              `json:"ndk"`          // Номер ДК
	NtoPhases     []NtoPhase       `json:"naptoph"`      // 1 - Массив привязки направлений к фазам (указаны те направления которые разрешены)
	TirToNaps     []TirToNap       `json:"tirtonap"`     // 2 - Привязка тиристоров к направлениям
	TminToPhases  []TimeToPhase    `json:"tmintophases"` // 3 - Привязка Тмин к номеру фазы
	RPUs          []RPU            `json:"rpus"`         // 4 - Привязка РПУ
	PromTaktBases []PromTakt       `json:"prombase"`     // 5 - Параметры промтакта для базовой РПУ
	PromTakt      []PromTakt       `json:"prom"`         // 6 - Параметры промтакта универсального
	NKs           [7]int           `json:"nks"`          // 8 - Привязка РПУ к дням недели
	CKs           []CK             `json:"cks"`          // 9 - Привязка РПУ к суточной карте
	Konf          []Mask           `json:"konf"`         //10 - конфликты для основного такта number-номер фазы
	TVPs          []TVP            `json:"tvps"`         //11 - Привязки ТВП
	TmaxToPhases  []TimeToPhase    `json:"tmaxtophases"` //13 - Привязка Тмах к номеру фазы
	Naps          map[int]bool     `json:"naps"`         // Состояние направлений
	StepPromtact  int              `json:"step"`         // Шаг в миллисекундах
	PromMake      PromMake         `json:"-"`            // Созданый промтакт переход
	PrMapBase     map[int]PromTakt `json:"-"`            // Промтакты базовые кешированные
	PrMapUn       map[int]PromTakt `json:"-"`            // Промтакты универсальные кешированные

}
type NtoPhase struct {
	NumPhase int   `json:"nph"`  // Номер фазы
	Naps     []int `json:"naps"` // Список направлений в фазе
}
type TirToNap struct {
	Number int   `json:"num"`    // Номер направления
	Type   int   `json:"type"`   // Тип направления 1-Транспортный 2-Поворотный постояный 3-Пешеходный 4-Трамвайный
	Green  int   `json:"green"`  // Номер зеленого тиристора
	Yellow int   `json:"yellow"` // Номер желтого тиристора
	Reds   []int `json:"reds"`   // Красный первый и так далее
}
type TimeToPhase struct {
	NumPhase int `json:"nphase"` // Номер фазы
	Tmin     int `json:"tmin"`   // Время фазы
}
type RPU struct {
	Number   int     `json:"number"` //Номер РПУ
	Tcycle   int     `json:"tcycle"` // Длительность цикла
	Continue bool    `json:"cont"`   // Признак продолжения цикла
	Phases   []Phase `json:"pahses"` // Фазы вызываемые в цикле
}
type Phase struct {
	Phase int `json:"phase"` // Номер  фазы
	Time  int `json:"time"`  // Время окончания фазы от начала цикла
}
type PromTakt struct {
	Nap        int `json:"nap"` //Номер направления
	GreenDop   int `json:"gd"`  //Начало зеленого дополнительного
	GreenBlink int `json:"gb"`  //Начало зеленого мигания
	Yellow     int `json:"yel"` //Начало желтого
	Red        int `json:"red"` //Начало красного
	RedYellow  int `json:"ry"`  //Начало красно-желтого
}
type Mask struct {
	Number int    `json:"number"` //Номер строки
	Mask   []byte `json:"mask"`   //Битовая маска
}
type CK struct {
	Number int    `json:"number"` //Номер карты
	Lines  []Line `json:"lines"`  //Переключения в суточной карте
}
type Line struct {
	Time   int `json:"time"`   //Время переключения
	Number int `json:"number"` //Номер РПУ
}
type TVP struct {
	Number int   `json:"number"` //Номер ТВП
	Wait   int   `json:"wait"`   //Номер тиристора к которому подключен "Ждите"
	Phases []int `json:"phases"` //Номера фаз в которые входит направление от ТВП

}
type KonfProm struct {
	NumPhase int   `json:"nph"`  // Номер фазы
	Naps     []int `json:"naps"` // Список направлений в фазе
}
