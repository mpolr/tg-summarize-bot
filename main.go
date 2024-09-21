package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Загружаем конфигурацию
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Создаем нового бота с токеном из конфигурации
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Инициализация канала для обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Основной цикл обработки обновлений
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Обработка команд
		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот для анализа чатов. Используй команду /вкратце для получения пересказа сообщений.")
			bot.Send(msg)
		case "help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используй команду /vkratse с опциональной датой (в формате YYYY-MM-DD), чтобы получить пересказ сообщений за выбранный день или последние сутки.")
			bot.Send(msg)
		case "vkratse":
			// Проверяем, указана ли дата
			args := update.Message.CommandArguments()
			var date time.Time
			if args != "" {
				date, err = time.Parse("2006-01-02", args)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажи корректную дату в формате YYYY-MM-DD.")
					bot.Send(msg)
					continue
				}
			} else {
				// Если дата не указана, получаем последние 24 часа
				date = time.Now().Add(-24 * time.Hour)
			}

			// Получаем сообщения за указанный период
			messages, err := getMessagesFromChat(bot, update.Message.Chat.ID, date)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении сообщений.")
				bot.Send(msg)
				continue
			}

			// Фильтруем сообщения по дате
			filteredMessages := getMessagesByDate(date, messages)
			if len(filteredMessages) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не найдено сообщений за указанную дату.")
				bot.Send(msg)
				continue
			}

			// Отправляем запрос на API Ollama для пересказа и получения тем
			summary, err := summarizeMessages(config.OllamaAPIURL, filteredMessages)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при обработке сообщений.")
				bot.Send(msg)
				continue
			}

			// Отправляем результат в чат
			dateStr := date.Format("2006-01-02")
			reply := fmt.Sprintf("Пересказ за %s:\n%s", dateStr, summary)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

// Фильтруем сообщения по дате
func getMessagesByDate(date time.Time, messages []MessageData) []MessageData {
	var results []MessageData
	for _, msg := range messages {
		if msg.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			results = append(results, msg)
		}
	}
	return results
}
