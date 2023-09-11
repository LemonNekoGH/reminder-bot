package config

import (
	"os"

	"github.com/joho/godotenv"
)

var Version = "0.1.0"

const (
	EnvDatabaseConnectUri = "REMINDER_BOT_DATABASE_URI"
	EnvTelegramBotToken   = "REMINDER_BOT_TOKEN" //nolint:gosec
)

type Config struct {
	DatabaseConnectUri string
	TelegramBotToken   string
}

func getEnv(key string) string {
	envMap, err := godotenv.Read()
	if err != nil {
		println("failed to read .env file")
	}

	value := envMap[key]
	if value != "" {
		return value
	}

	return os.Getenv(key)
}

func NewConfig() func() *Config {
	return func() *Config {
		return &Config{
			DatabaseConnectUri: getEnv(EnvDatabaseConnectUri),
			TelegramBotToken:   getEnv(EnvTelegramBotToken),
		}
	}
}
