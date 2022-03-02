package bin

type CMK struct {
	BaseOption    int           `json:"bopt"`         // Базовые опции
	BaseRPU       int           `json:"brpu"`         // Номер базовой РПУ
	NumTU         int           `json:"ntu"`          // Номер ТУ
	NumDK         int           `json:"ndk"`          // Номер ДК
	NtoPhases     []NtoPhase    `json:"naptoph"`      // 1 - Массив привязки направлений к фазам
	TirToNaps     []TirToNap    `json:"tirtonap"`     // 2 - Привязка тиристоров к направлениям
	TminToPhases  []TimeToPhase `json:"tmintophases"` // 3 - Привязка Тмин к номеру фазы
	RPUs          []RPU         `json:"rpus"`         // 4 - Привязка РПУ
	PromTaktBases []PromTakt    `json:"prombase"`     // 5 - Параметры промтакта для базовой РПУ
	PromTakt      []PromTakt    `json:"prom"`         // 6 - Параметры промтакта
	RedGroups     []Mask        `json:"redgs"`        // 7 - Массив групп красных ламп number-номер группы
	NKs           [7]int        `json:"nks"`          // 8 - Привязка РПУ к дням недели
	CKs           []CK          `json:"cks"`          // 9 - Привязка РПУ к суточной карте
	Konf          []Mask        `json:"konf"`         //10 - конфликты для основного такта number-номер фазы
	TVPs          []TVP         `json:"tvps"`         //11 - Привязки ТВП
	GreenBlinks   []NtoPhase    `json:"gbs"`          //12 - Привязка направлений зеленого мигания
	TmaxToPhases  []TimeToPhase `json:"tmaxtophases"` //13 - Привязка Тмах к номеру фазы
	GreenToOut    []Mask        `json:"greenout"`     //14 - привязки зеленых ламп number-1
	PowerToOut    []Mask        `json:"powerout"`     //15 - привязки тиристоров по мощности number-1 1-тиристор контроль частичное перегорание
	ReedToOut     []Mask        `json:"redout"`       //16 - привязки красных ламп к тиристорам number-1
	YellowToOut   []Mask        `json:"yellowout"`    //17 - привязки желтых ламп к тиристорам number-1
	Phases        []int         `json:"phases"`       //18 - номера фаз базовой РПУ
	Switches      []int         `json:"switches"`     //19 - Массив возможных переходов ??
	KonfProms     []KonfProm    `json:"konfprom"`     //20 - Матрица конфликтных  для пром такта
	Brokens       []Mask        `json:"brokens"`      //22 - лампы контроллируемые на полное перегорание
}
type NtoPhase struct {
	NumPhase int   `json:"nph"`  // Номер фазы
	Naps     []int `json:"naps"` // Список направлений в фазе
}
type TirToNap struct {
	Number int   `json:"num"`    // Номер направления
	Green  int   `json:"green"`  // Номер зеленого тиристора
	Yellow int   `json:"yellow"` // Номер желтого тиристора
	Reds   []int `json:"reds"`   // Красный первый и так далее
}
type TimeToPhase struct {
	NumPhase int `json:"nphase"` // Номер фазы
	Tmin     int `json:"tmin"`   // Время фазы
}
type RPU struct {
	Number int `json:"number"` //Номер РПУ
	Tcycle int `json:"tcycle"` // Длительность цикла
	// Del1     int     `json:"d1"`     //Остаток от деленения 256 на время цикла
	// Del2     int     `json:"d2"`     //Остаток от деленения 65536 на время цикла
	Continue bool    `json:"cont"`   // Признак продолжения цикла
	Phases   []Phase `json:"pahses"` // Фазы вызываемые в цикле
}
type Phase struct {
	TVP1  int `json:"tvp1"` // Номер фазы с ТВП1
	TVP2  int `json:"tvp2"` // Номер фазы с ТВП2
	ZamPh int `json:"zam"`  // Номер замещающей фазы
	Time  int `json:"time"` // Время окончания фазы от начала цикла
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
