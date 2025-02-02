package config

import (
    "log"
    "os"
)

var (
    // Убираем ненужную переменную
    DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
)

type Config struct {
    BotToken string
}

func GetConfig() *Config {
    // Используем переменную окружения TELEGRAM_BOT_TOKEN
    botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
    if botToken == "" {
        log.Fatal("TELEGRAM_BOT_TOKEN не найден в переменных окружения")
    }

    return &Config{
        BotToken: botToken,
    }
}

func InitConfig() {
    // Проверка на обязательные переменные окружения
    if DBConnectionString == "" {
        log.Fatal("Missing environment variable: DB_CONNECTION_STRING")
    }
}
