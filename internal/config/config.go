package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"

	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env"  env-default:"development"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"internal/db"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address"  env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout"  env-default:"4s"`
	Debug       bool          `yaml:"debug"  env-default:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"10s"`
}

func InitConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist:%s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
