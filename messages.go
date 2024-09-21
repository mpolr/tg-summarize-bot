package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageData хранит информацию о сообщении
type MessageData struct {
	ID               int
	Date             time.Time
	Text             string
	SenderUserName   string
	SenderUserID     int64
	IsReply          bool
	ReplyToMessageID int
}

// getTelegramUserName возвращает имя отправителя сообщения
func getTelegramUserName(sender *tgbotapi.User) string {
	if sender == nil {
		return "<unknown>"
	}
	if sender.FirstName != "" && sender.LastName != "" {
		return sender.FirstName + " " + sender.LastName
	} else if sender.FirstName != "" {
		return sender.FirstName
	} else if sender.LastName != "" {
		return sender.LastName
	}
	return "<unknown>"
}

// getMessagesFromChat получает сообщения из чата за указанный период времени
func getMessagesFromChat(bot *tgbotapi.BotAPI, chatID int64, date time.Time) ([]MessageData, error) {
	var history []MessageData

	// Получаем обновления из чата
	updates, err := bot.GetUpdates(tgbotapi.UpdateConfig{})
	if err != nil {
		return nil, err
	}

	for _, update := range updates {
		if update.Message != nil && update.Message.Chat.ID == chatID {
			messageDate := time.Unix(int64(update.Message.Date), 0)
			if messageDate.After(date) {
				sender := update.Message.From
				data := MessageData{
					ID:             update.Message.MessageID,
					Date:           messageDate,
					Text:           update.Message.Text,
					SenderUserName: getTelegramUserName(sender),
					SenderUserID:   sender.ID,
					IsReply:        update.Message.ReplyToMessage != nil,
				}
				if data.IsReply {
					data.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
				}
				history = append(history, data)
			}
		}
	}

	// Реверсируем историю, чтобы получить порядок от старых к новым
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return history, nil
}
