package repository

import (
	"database/sql"
	"telegram/models"
	"log"
)

// UserRepository — структура для работы с пользователями
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository — конструктор для UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser — создаем нового пользователя в базе данных
func (r *UserRepository) CreateUser(user models.User) (models.User, error) {
	var newUser models.User
	query := `INSERT INTO users (id, telegram_username, telegram_first_name, telegram_last_name) 
	          VALUES ($1, $2, $3, $4) RETURNING id, telegram_username, telegram_first_name, telegram_last_name`
	err := r.DB.QueryRow(query, user.ID, user.TelegramUsername, user.TelegramFirstName, user.TelegramLastName).Scan(
		&newUser.ID, &newUser.TelegramUsername, &newUser.TelegramFirstName, &newUser.TelegramLastName)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return newUser, err
	}
	return newUser, nil
}

// GetUserByID — ищем пользователя по ID
func (r *UserRepository) GetUserByID(userID int64) (models.User, error) {
	var user models.User
	query := `SELECT id, telegram_username, telegram_first_name, telegram_last_name 
	          FROM users WHERE id = $1`
	err := r.DB.QueryRow(query, userID).Scan(&user.ID, &user.TelegramUsername, &user.TelegramFirstName, &user.TelegramLastName)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return user, err
	}
	return user, nil
}
