package models

// User представляет пользователя бота
type User struct {
	ID                int64  `json:"id"`
	TelegramUsername  string `json:"telegram_username"`
	TelegramFirstName string `json:"telegram_first_name"`
	TelegramLastName  string `json:"telegram_last_name"`
}
