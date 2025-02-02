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
func (repo *UserRepository) CreateUser(user models.User) (*models.User, error) {
	// Проверяем, существует ли уже пользователь с таким Telegram ID
	var existingUser models.User
	err := repo.DB.QueryRow("SELECT id, telegram_username, telegram_first_name, telegram_last_name FROM users WHERE id = $1", user.ID).Scan(&existingUser.ID, &existingUser.TelegramUsername, &existingUser.TelegramFirstName, &existingUser.TelegramLastName)
	
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, err
	}

	// Если пользователь уже существует, возвращаем его
	if existingUser.ID != 0 {
		return &existingUser, nil
	}

	// Создаем нового пользователя, если его нет в базе
	_, err = repo.DB.Exec("INSERT INTO users (id, telegram_username, telegram_first_name, telegram_last_name) VALUES ($1, $2, $3, $4)", user.ID, user.TelegramUsername, user.TelegramFirstName, user.TelegramLastName)
	if err != nil {
		return nil, err
	}

	return &user, nil
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
