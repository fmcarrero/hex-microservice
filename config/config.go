package config

type Configuration struct {
	Database struct {
		Host string `yaml:"host", envconfig:"DATABASE_HOST"`
		Port string `yaml:"port", envconfig:"DATABASE_PORT"`
	} `yaml:"database"`
	Server struct {
		Port string `yaml:"port", envconfig:"SERVER_PORT"`
	} `yaml:"server"`
}
