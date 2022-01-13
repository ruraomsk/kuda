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
}
type SetupBrams struct {
	DbPath string `toml:"dbpath"`
	Step   int    `toml:"step"`
}
type WatchDog struct {
	Step int `toml:"step"`
}
type Hardware struct {
	Step int `toml:"step"`
}
type Netware struct {
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

func init() {
}
