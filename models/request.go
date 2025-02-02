package models

// Request представляет заявку пользователя
type Request struct {
	ID          int64    `json:"id"`
	UserID      int64  `json:"user_id"`
	Service string `json:"service"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
