package setup

var (
	Set *Setup
)

type Setup struct {
	LogPath    string     `toml:"logpath"`
	SetupBrams SetupBrams `toml:"setBrams"`
	WatchDog   WatchDog   `toml:"watchdog"`
	Hardware   Hardware   `toml:"hardware"`
	Netware    Netware    `toml:"netware"`
	Vpu        Vpu        `toml:"vpu"`
	Counter    Counter    `toml:"counter"`
}
type SetupBrams struct {
	DbPath string `toml:"dbpath"`
	Step   int    `toml:"step"`
}
type WatchDog struct {
	Step int `toml:"step"`
}
type Hardware struct {
	Step    int    `toml:"step"`
	Connect string `toml:"connect"`
	SPort   int    `toml:"sport"`
	C8count int    `toml:"count"`
	LongKK  int    `toml:"longkk"`
	PinOS   int    `toml:"pinos"`
	PinYB   int    `toml:"pinyb"`
}
type Netware struct {
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}
type Vpu struct {
	Step    int    `toml:"step"`
	Connect string `toml:"connect"`
	SPort   int    `toml:"sport"`
}
type Counter struct {
	Step    int    `toml:"step"`
	Connect string `toml:"connect"`
	SPort   int    `toml:"sport"`
}

func init() {
}
