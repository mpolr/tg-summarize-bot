package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Форматируем ответ для отправки в чат
func formatResponse(userName, summary string) string {
	return "@" + userName + ", вот краткий пересказ: " + summary
}

// Текст команды помощи
func getHelpText() string {
	return "Я бот, который пересказывает сообщения за указанную дату. Используйте @бот и дату в формате YYYY-MM-DD."
}

// Функция для проверки, упомянут ли бот в сообщении
func isBotMentioned(message *tgbotapi.Message, bot *tgbotapi.BotAPI) bool {
	if message.Entities == nil {
		return false
	}

	// Проходим по сущностям сообщения (например, упоминания)
	for _, entity := range message.Entities {
		if entity.Type == "mention" {
			// Проверяем, совпадает ли упоминание с именем бота
			mention := message.Text[entity.Offset : entity.Offset+entity.Length]
			if mention == "@"+bot.Self.UserName {
				return true
			}
		}
	}

	return false
}
