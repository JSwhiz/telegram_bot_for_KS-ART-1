package main

import (
	"log"
	"os"
	"telegram/handlers"
	"telegram/repository"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"database/sql"
	_ "github.com/lib/pq"
)

func main() {
	// Получаем токен
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не установлен")
	}

	// Подключаемся к базе данных
	connStr := os.Getenv("DB_CONNECTION_STRING")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем репозитории для работы с пользователями и заявками
	userRepository := repository.NewUserRepository(db)
	requestRepository := repository.NewRequestRepository(db)

	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Получаем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.Text == "/start" {
				sendWelcomeMessage(bot, update.Message.Chat.ID)
			}
		}

		// Обрабатываем callback-запросы
		if update.CallbackQuery != nil {
			handlers.HandleServiceSelection(bot, update, userRepository, requestRepository)
		}
	}
}

// sendWelcomeMessage отправляет приветственное сообщение с файлами и кнопками
func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64) {
	// Текст приветствия с описанием и ссылкой на PDF
	messageText := "👋 Привет! Мы — KS-ART, начинающий IT-стартап, который создает технологические решения для бизнеса.\n\n" +
		"📄 Мы прикрепили PDF-документ, который расскажет больше о нашей компании.\n\n" +
		"💼 Выберите услугу, с которой хотите работать:\n\n" 

	// Отправляем PDF файл с текстом
	pdfFile := tgbotapi.NewDocument(chatID, tgbotapi.FilePath("../logs/Почему именно KS-ART?.pdf"))
	pdfFile.Caption = messageText

	// Кнопки с выбором услуг (вертикально)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌍 Веб-разработка", "service_web"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 Разработка приложений", "service_app"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎨 Дизайн", "service_design"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🖥 Поставка оборудования", "service_hardware"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔧 IT сопровождение", "service_it"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Другое", "service_other"),
		),
	)

	// Добавляем клавиатуру к сообщению
	pdfFile.ReplyMarkup = keyboard

	// Отправляем документ
	bot.Send(pdfFile)
}
