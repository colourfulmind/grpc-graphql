package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	postgres "ozon/internal/storage/postgres"
	"time"
)

type Config struct {
	Env      string            `yaml:"env"`
	TokenTTL time.Duration     `yaml:"token_ttl"`
	GRPC     GRPCConfig        `yaml:"grpc"`
	Postgres postgres.Postgres `yaml:"postgres"`
	Clients  Clients           `yaml:"clients"`
}

type GRPCConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Clients struct {
	GRPCClient GRPCClient `yaml:"grpc"`
}

type GRPCClient struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	RetriesCount int           `yaml:"retries_count"`
	Timeout      time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := FetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists: " + err.Error())
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

func FetchConfigPath() (path string) {
	flag.StringVar(&path, "config", "", "path to config")
	flag.Parse()
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}
