package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env" env-default:"local"`
	GRPC GRPC
}

type GRPC struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() Config {
	path := fetchPathConfig()
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	if path == "" {
		panic("CONFIG_PATH is empty")
	}
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("err to load config fail")
	}
	return cfg
}

func fetchPathConfig() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		panic("config file path is empty")
	}
	return path
}
