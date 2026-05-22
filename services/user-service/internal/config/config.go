package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	TokenTTL time.Duration  `yaml:"token_ttl" env-required:"true"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Postgres PostgresConfig `yaml:"postgres"`
	Infra    InfraConfig    `yaml:"infra"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type PostgresConfig struct {
	Host     string        `yaml:"host"     env:"POSTGRES_HOST"`
	Port     string        `yaml:"port"     env-default:"5432"`
	User     string        `yaml:"user"     env-required:"true"`
	Password string        `yaml:"password" env-required:"true"`
	Database string        `yaml:"database" env-required:"true"`
	Timeout  time.Duration `yaml:"timeout"  env-required:"true"`
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
