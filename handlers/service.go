package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"telegram/models"
	"telegram/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик выбора услуги
func HandleServiceSelection(bot *tgbotapi.BotAPI, update tgbotapi.Update, userRepository *repository.UserRepository, requestRepository *repository.RequestRepository) {
	log.Printf("Получен запрос на выбор услуги от пользователя %d", update.CallbackQuery.From.ID)
	user := models.User{
		ID:                update.CallbackQuery.From.ID,
		TelegramUsername:  update.CallbackQuery.From.UserName,
		TelegramFirstName: update.CallbackQuery.From.FirstName,
		TelegramLastName:  update.CallbackQuery.From.LastName,
	}

	// Сохраняем пользователя в базе данных (если его еще нет)
	_, err := userRepository.CreateUser(user)
	if err != nil {
		log.Printf("Ошибка при создании пользователя %d: %v", user.ID, err)
		return
	}

	// Получаем выбранную услугу из данных кнопки
	service := update.CallbackQuery.Data
	log.Printf("Пользователь %d выбрал услугу: %s", user.ID, service)

	// Создаем заявку в базе данных
	request, err := requestRepository.CreateRequest(user.ID, service)
	if err != nil {
		log.Printf("Ошибка при создании заявки для пользователя %d: %v", user.ID, err)
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
	log.Printf("Отправка запроса на описание для заявки с ID %d пользователю %d", requestID, chatID)
	messageText := fmt.Sprintf("📄 Опишите вашу задачу для заявки №%d.", requestID)
	msg := tgbotapi.NewMessage(chatID, messageText)
	bot.Send(msg)
}

// Пример метода для получения ID запроса для пользователя
func getRequestIDForUser(userID int64, requestRepository *repository.RequestRepository) (int64, error) {
	log.Printf("Получение ID запроса для пользователя %d", userID)
	var requestID int64
	log.Printf("Запрос к БД: userID=%d, статус='new'", userID)
	query := `SELECT id FROM requests WHERE user_id = $1 AND status = 'new' LIMIT 1`
	err := requestRepository.DB.QueryRow(query, userID).Scan(&requestID)
	if err != nil {
		log.Printf("Ошибка при получении ID запроса для пользователя %d: %v", userID, err)
		return 0, fmt.Errorf("failed to get request ID for user %d: %w", userID, err)
	}
	log.Printf("ID запроса для пользователя %d: %d", userID, requestID)
	return requestID, nil
}

func HandleDescription(bot *tgbotapi.BotAPI, update tgbotapi.Update, requestRepository *repository.RequestRepository) {
	
	description := update.Message.Text
	userID := update.Message.From.ID
	log.Printf("HandleDescription вызван для userID=%d с текстом: %s", userID, update.Message.Text)
	log.Printf("Получено описание от пользователя %d: %s", userID, description)

	// Декодируем Unicode-символы в читаемый текст
	decodedDescription := decodeUnicode(description)

	// Логируем декодированное описание
	log.Printf("Декодированное описание: %s", decodedDescription)

	// Получаем ID заявки для текущего пользователя
	requestID, err := getRequestIDForUser(userID, requestRepository)
	if err != nil {
		log.Printf("Ошибка получения ID заявки для пользователя %d: %v", userID, err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при получении вашего запроса.")
		bot.Send(msg)
		return
	}

	// Если описание пустое, просим ввести его
	if decodedDescription == "" {
		log.Printf("Описание пустое, просим пользователя %d ввести описание.", userID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите описание для вашей заявки.")
		bot.Send(msg)
		return
	}

	// Обновляем описание заявки в базе данных
	err = requestRepository.UpdateRequestDescription(requestID, decodedDescription)
	if err != nil {
		log.Printf("Ошибка при обновлении описания заявки с ID %d: %v", requestID, err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при сохранении описания.")
		bot.Send(msg)
		return
	}

	// Логируем успешное обновление
	log.Printf("Описание заявки с ID %d успешно обновлено для пользователя %d", requestID, userID)

	// Отправляем сообщение об успешном сохранении
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваше описание успешно сохранено. Спасибо!")
	bot.Send(msg)
}

// decodeUnicode — декодирует строку с Unicode-символами в читаемый текст
func decodeUnicode(input string) string {
	var decoded string
	err := json.Unmarshal([]byte(`"`+input+`"`), &decoded)
	if err != nil {
		log.Printf("Ошибка при декодировании строки: %v", err)
		return input // Возвращаем оригинальную строку, если ошибка декодирования
	}
	return decoded
}
