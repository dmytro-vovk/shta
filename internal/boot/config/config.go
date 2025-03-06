package config

import (
	"os"

	"github.com/dmytro-vovk/envset"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database  Database  `yaml:"database"`
	Cache     Cache     `yaml:"cache"`
	WebServer WebServer `yaml:"webserver"`
}

type Database struct {
	Host     string `default:"localhost"            env:"DB_HOST"     yaml:"host"`
	Port     string `default:"5432"                 env:"DB_PORT"     yaml:"port"`
	User     string `default:"appuser"              env:"DB_USERNAME" yaml:"user"`
	Password string `default:"sEcRetPaSs8371238642" env:"DB_PASSWORD" yaml:"pass"`
	Name     string `default:"appdb"                env:"DB_NAME"     yaml:"name"`
}

type Cache struct{}

type WebServer struct {
	Listen string `default:"0.0.0.0:8080" env:"LISTEN" yaml:"listen"`
}

func Load(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := envset.Set(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
