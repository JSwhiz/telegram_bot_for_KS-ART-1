package config

import (
    "log"
    "os"
)

var (
    TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

    DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
)

func InitConfig() {
    if TelegramBotToken == "" || DBConnectionString == "" {
        log.Fatal("Missing environment variables: TELEGRAM_BOT_TOKEN or DB_CONNECTION_STRING")
    }
}
