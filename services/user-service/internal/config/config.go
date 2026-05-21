package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env       string         `yaml:"env" env-default:"local"`
	TokenTTL  time.Duration  `yaml:"token_ttl" env-required:"true"`
	JWTSecret string         `env:"JWT_SECRET" env-required:"true"`
	GRPC      GRPCConfig     `yaml:"grpc"`
	Postgres  PostgresConfig `yaml:"postgres"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type PostgresConfig struct {
	Host     string        `yaml:"host"     env-required:"true"`
	Port     string        `yaml:"port"     env-default:"5432"`
	User     string        `yaml:"user"     env-required:"true"`
	Password string        `env:"POSTGRES_PASSWORD" env-required:"true"`
	Database string        `yaml:"database" env-required:"true"`
	Timeout  time.Duration `yaml:"timeout"  env-required:"true"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	path := os.Getenv("CONFIG_PATH")
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
