package config

type Server struct {
	Addr    string  `yaml:"addr"`
	Service Service `yaml:"service"`
	Logger  Logger  `yaml:"logger"`
}

type Service struct {
	Dir string `yaml:"dir"`
}
