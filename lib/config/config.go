package config

type Config struct {
	Interval int
	Port     int
	TimeOut  int
	Targets  []TargetConfig
}

type TargetConfig struct {
	Url  string
	Name string
}
