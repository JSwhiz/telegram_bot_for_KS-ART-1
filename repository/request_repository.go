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
	return &RequestRepository{DB: db}
}

// CreateRequest — создаем новую заявку
func (r *RequestRepository) CreateRequest(userID int64, service string) (models.Request, error) {
	var request models.Request
	query := `INSERT INTO requests (user_id, service, status) VALUES ($1, $2, $3) RETURNING id, user_id, service, status`
	err := r.DB.QueryRow(query, userID, service, "new").Scan(&request.ID, &request.UserID, &request.Service, &request.Status)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return request, err
	}
	return request, nil
}

// GetRequestsByStatus — получает заявки по статусу
func (r *RequestRepository) GetRequestsByStatus(status string) ([]models.Request, error) {
	rows, err := r.DB.Query("SELECT id, user_id, service, description, status FROM requests WHERE status = $1", status)
	if err != nil {
		log.Printf("Error fetching requests by status: %v", err)
		return nil, err
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var request models.Request
		if err := rows.Scan(&request.ID, &request.UserID, &request.Service, &request.Description, &request.Status); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error with rows iteration: %v", err)
		return nil, err
	}

	return requests, nil
}


// GetRequestsWithoutDescription — получает заявки без описания
func (r *RequestRepository) GetRequestsWithoutDescription() ([]models.Request, error) {
	rows, err := r.DB.Query("SELECT id, user_id, service, description, status FROM requests WHERE description IS NULL OR description = ''")
	if err != nil {
		log.Printf("Error fetching requests without description: %v", err)
		return nil, err
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var request models.Request
		if err := rows.Scan(&request.ID, &request.UserID, &request.Service, &request.Description, &request.Status); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error with rows iteration: %v", err)
		return nil, err
	}

	return requests, nil
}

// UpdateRequestStatus — обновляет статус заявки
func (r *RequestRepository) UpdateRequestStatus(requestID int64, status string) error {
	_, err := r.DB.Exec("UPDATE requests SET status = $1 WHERE id = $2", status, requestID)
	if err != nil {
		log.Printf("Error updating request status: %v", err)
		return err
	}
	return nil
}

// Метод для обновления описания заявки в базе данных
func (r *RequestRepository) UpdateRequestDescription(requestID int64, description string) error {
	query := `UPDATE requests SET description = $1 WHERE id = $2`
	_, err := r.db.Exec(query, description, requestID)
	if err != nil {
		return fmt.Errorf("failed to update description: %w", err)
	}
	return nil
}


