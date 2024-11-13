package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Env         string     `json:"env" env:"ENV" env-default:"local"` // local, dev, prod
	StoragePath string     `json:"storage_path" env:"STORAGE_PATH" env-default:"/storage/sqlite/nms.db"`
	HTTPServer  HTTPServer `json:"http_server"`
}

type HTTPServer struct {
	Address string `json:"address" env:"ADDRESS" env-default:"localhost:8080"`
	// Timeout     time.Duration `json:"timeout" env:"TIMEOUT" env-default:"4s"`
	// IdleTimeout time.Duration `json:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	defaultConfigPath := "config/local.json"

	if configPath == "" {
		fmt.Println("CONFIG_PATH is not set, using default config:", defaultConfigPath)
		configPath = defaultConfigPath
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file at %s does not exist", configPath)
	}

	var cfg Config

	jsonFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	if err := json.Unmarshal(jsonFile, &cfg); err != nil {
		log.Fatalf("failed to parse config at %s: %v", configPath, err)
	}

	return &cfg
}
