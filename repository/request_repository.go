package repository

import (
	"database/sql"
	"log"
	"telegram/models"
	"fmt"
)

// RequestRepository — структура для работы с заявками
type RequestRepository struct {
	DB *sql.DB
}

// NewRequestRepository — конструктор для RequestRepository
func NewRequestRepository(db *sql.DB) *RequestRepository {
	log.Println("Создание нового репозитория заявок с базой данных.")
	return &RequestRepository{DB: db}
}

// CreateRequest — создаем новую заявку
func (r *RequestRepository) CreateRequest(userID int64, service string) (models.Request, error) {
	var request models.Request
	log.Printf("Создание новой заявки для пользователя %d, услуга: %s", userID, service)
	query := `INSERT INTO requests (user_id, service, status) VALUES ($1, $2, $3) RETURNING id, user_id, service, status`
	err := r.DB.QueryRow(query, userID, service, "new").Scan(&request.ID, &request.UserID, &request.Service, &request.Status)
	if err != nil {
		log.Printf("Ошибка при создании заявки для пользователя %d: %v", userID, err)
		return request, err
	}
	log.Printf("Заявка успешно создана с ID %d для пользователя %d", request.ID, userID)
	return request, nil
}

// GetRequestsByStatus — получает заявки по статусу
func (r *RequestRepository) GetRequestsByStatus(status string) ([]models.Request, error) {
	log.Printf("Запрос заявок с статусом: %s", status)
	rows, err := r.DB.Query("SELECT id, user_id, service, description, status FROM requests WHERE status = $1", status)
	if err != nil {
		log.Printf("Ошибка при получении заявок по статусу %s: %v", status, err)
		return nil, err
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var request models.Request
		if err := rows.Scan(&request.ID, &request.UserID, &request.Service, &request.Description, &request.Status); err != nil {
			log.Printf("Ошибка при сканировании строки с заявкой: %v", err)
			return nil, err
		}
		requests = append(requests, request)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return nil, err
	}
	log.Printf("Получено %d заявок с статусом %s", len(requests), status)
	return requests, nil
}

// GetRequestsWithoutDescription — получает заявки без описания
func (r *RequestRepository) GetRequestsWithoutDescription() ([]models.Request, error) {
	log.Println("Запрос заявок без описания.")
	rows, err := r.DB.Query("SELECT id, user_id, service, description, status FROM requests WHERE description IS NULL OR description = ''")
	if err != nil {
		log.Printf("Ошибка при получении заявок без описания: %v", err)
		return nil, err
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var request models.Request
		if err := rows.Scan(&request.ID, &request.UserID, &request.Service, &request.Description, &request.Status); err != nil {
			log.Printf("Ошибка при сканировании строки с заявкой: %v", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return nil, err
	}

	log.Printf("Получено %d заявок без описания.", len(requests))
	return requests, nil
}

// UpdateRequestStatus — обновляет статус заявки
func (r *RequestRepository) UpdateRequestStatus(requestID int64, status string) error {
	log.Printf("Обновление статуса заявки с ID %d на %s", requestID, status)
	_, err := r.DB.Exec("UPDATE requests SET status = $1 WHERE id = $2", status, requestID)
	if err != nil {
		log.Printf("Ошибка при обновлении статуса заявки с ID %d: %v", requestID, err)
		return err
	}
	log.Printf("Статус заявки с ID %d успешно обновлен на %s", requestID, status)
	return nil
}

// Метод для обновления описания заявки в базе данных
func (r *RequestRepository) UpdateRequestDescription(requestID int64, description string) error {
	log.Printf("Обновление описания заявки с ID %d: %s", requestID, description)
	query := `UPDATE requests SET description = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, description, requestID)
	if err != nil {
		log.Printf("Ошибка при обновлении описания заявки с ID %d: %v", requestID, err)
		return fmt.Errorf("не удалось обновить описание заявки: %w", err)
	}
	log.Printf("Описание заявки с ID %d успешно обновлено.", requestID)
	log.Printf("Обновляем описание заявки %d: %s", requestID, description)
	return nil
}
