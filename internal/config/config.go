package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env           string `yaml:"env" env:"ENV" envDefault:"local"`
	StoragePath   string `yaml:"storage_path" env-required:"true"`
	TgBotApiHost  string `yaml:"tg_bot_api_host" env-required:"true"`
	TgBotApiToken string `yaml:"tg_bot_api_token" env-required:"true" env:"TG_BOT_API_TOKEN"`
	BatchSize     int    `yaml:"batch_size" envDefault:"100"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return &cfg
}
