package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string            `yaml:"env" env-default:"local"`
	HTTP        HTTPConfig        `yaml:"http"`
	UserService UserServiceConfig `yaml:"user_service"`
	Infra       InfraConfig       `yaml:"infra"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type UserServiceConfig struct {
	Addr    string        `yaml:"addr"    env:"USER_SERVICE_ADDR"`
	Timeout time.Duration `yaml:"timeout"`
}

type InfraConfig struct {
	JWTSecret string `yaml:"jwt_secret" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
