package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type JWTTokenConfig struct {
	Secret             string        `yaml:"secret" json:"secret" env:"SECRET"`
	ExpirationDuration time.Duration `yaml:"expiration_duration" json:"expiration_duration" env:"EXPIRATION_DURATION"`
}

type Config struct {
	Env            string          `yaml:"env" json:"env" env:"ENV" env-default:"local"` // local, dev, prod
	StoragePath    string          `yaml:"storage_path" json:"storage_path" env:"STORAGE_PATH" env-default:"/storage/sqlite/storage.db"`
	HTTPServer     HTTPServer      `yaml:"http_server" json:"http_server"`
	JWTTokenConfig *JWTTokenConfig `yaml:"jwt_token" json:"jwt_token"`
}

type HTTPServer struct {
	Address     string        `json:"address" yaml:"address" env:"adress" env-default:"localhost:8080"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout" env:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `json:"idle_timeout" yaml:"idle_timeout" env:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	defaultConfigPath := "configs/local.yml"

	if configPath == "" {
		fmt.Println("CONFIG_PATH is not set, using default config:", defaultConfigPath)
		configPath = defaultConfigPath
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file at %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to parse config at %s, err: %v", configPath, err)
	}

	return &cfg
}
