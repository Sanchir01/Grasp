package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"log/slog"

	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env"  env-default:"development"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"internal/db"`
	HttpServer  `yaml:"http_server"`
	Errors      `yaml:"errors"`
}

type HttpServer struct {
	Timeout     time.Duration `yaml:"timeout"  env-default:"4s"`
	Host        string        `yaml:"host"  env-default:"localhost"`
	Port        string        `yaml:"port"  env-default:"8000"`
	Debug       bool          `yaml:"debug"  env-default:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"60s"`
}

type Errors struct {
	Unauthorized ErrorsBody `yaml:"unauthorized"  env-default:"Unauthorized"`
	NotFound     ErrorsBody `yaml:"not_found"  env-default:"Not found"`
}

type ErrorsBody struct {
	message string `yaml:"message"`
	code    int    `yaml:"code"`
}

func InitConfig() *Config {
	if err := godotenv.Load(".development.env"); err != nil {
		slog.Error("ошибка при инициализации переменных окружения", err.Error())
	}
	configPath := os.Getenv("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist:%s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
