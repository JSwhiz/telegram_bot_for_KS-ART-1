package handlers

import (
	"fmt"
	"telegram/models"
	"telegram/repository"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик выбора услуги
func HandleServiceSelection(bot *tgbotapi.BotAPI, update tgbotapi.Update, userRepository *repository.UserRepository, requestRepository *repository.RequestRepository) {
	// Получаем данные пользователя из обновления
	user := models.User{
		ID:                update.CallbackQuery.From.ID,
		TelegramUsername:  update.CallbackQuery.From.UserName,
		TelegramFirstName: update.CallbackQuery.From.FirstName,
		TelegramLastName:  update.CallbackQuery.From.LastName,
	}

	// Сохраняем пользователя в базе данных (если его еще нет)
	_, err := userRepository.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}

	// Получаем выбранную услугу из данных кнопки
	service := update.CallbackQuery.Data

	// Создаем заявку в базе данных
	request, err := requestRepository.CreateRequest(user.ID, service)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	// Отправляем сообщение о подтверждении
	messageText := fmt.Sprintf("Вы выбрали услугу: *%s*\n\nПожалуйста, опишите вашу задачу.", service)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messageText)
	msg.ParseMode = "Markdown" // Для работы с Markdown
	bot.Send(msg)

	// Отправляем запрос на описание задачи
	sendDescriptionRequest(bot, update.CallbackQuery.Message.Chat.ID, request.ID)
}

// sendDescriptionRequest — отправляем запрос на описание задачи
func sendDescriptionRequest(bot *tgbotapi.BotAPI, chatID int64, requestID int64) {
	messageText := fmt.Sprintf("📄 Опишите вашу задачу для заявки №%d.", requestID)
	msg := tgbotapi.NewMessage(chatID, messageText)
	bot.Send(msg)
}

// Пример метода для получения ID запроса для пользователя
func getRequestIDForUser(userID int, requestRepository *repository.RequestRepository) (int64, error) {
	var requestID int64
	query := `SELECT id FROM requests WHERE user_id = $1 AND status = 'New' LIMIT 1`
	err := requestRepository.db.QueryRow(query, userID).Scan(&requestID)
	if err != nil {
		return 0, fmt.Errorf("failed to get request ID for user %d: %w", userID, err)
	}
	return requestID, nil
}


// Обработчик ввода описания заявки
func HandleDescriptionInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, requestRepository *repository.RequestRepository) {
	// Получаем текст сообщения (это описание заявки)
	description := update.Message.Text

	// Получаем ID заявки для текущего пользователя
	userID := update.Message.From.ID // Получаем ID пользователя Telegram
	requestID, err := getRequestIDForUser(userID, requestRepository) // Нужно добавить этот метод
	if err != nil {
		log.Printf("Error fetching request ID for user %d: %v", userID, err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при получении вашего запроса.")
		bot.Send(msg)
		return
	}

	// Если описание пустое, просим ввести его
	if description == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите описание для вашей заявки.")
		bot.Send(msg)
		return
	}

	// Обновляем описание заявки в базе данных
	err = requestRepository.UpdateRequestDescription(requestID, description)
	if err != nil {
		log.Printf("Error updating request description: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при сохранении описания.")
		bot.Send(msg)
		return
	}

	// Отправляем сообщение об успешном сохранении
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваше описание успешно сохранено. Спасибо!")
	bot.Send(msg)
}
