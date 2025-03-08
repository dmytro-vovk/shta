package config

import (
	"fmt"
	"os"
	"time"

	"github.com/dmytro-vovk/envset"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database  Database  `yaml:"database"`
	WebServer WebServer `yaml:"webserver"`
	Settings  struct {
		TopURLs          int           `default:"10"  yaml:"topUrls"`
		ConcurrencyLimit int           `default:"3"   yaml:"httpConcurrencyLimit"`
		VerifyEvery      time.Duration `default:"10s" yaml:"verifyEvery"`
	} `yaml:"settings"`
}

type Database struct {
	Host     string `default:"localhost" env:"DB_HOST"     yaml:"host"`
	Port     string `default:"5432"      env:"DB_PORT"     yaml:"port"`
	User     string `default:"appuser"   env:"DB_USERNAME" yaml:"user"`
	Password string `default:"-"         env:"DB_PASSWORD" yaml:"pass"`
	Name     string `default:"appdb"     env:"DB_NAME"     yaml:"name"`
}

func (db Database) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}

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
