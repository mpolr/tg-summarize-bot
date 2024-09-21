package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// OllamaRequest структура для формирования запроса к API Ollama
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Функция для отправки сообщений на REST API Ollama
func summarizeMessages(apiURL string, messages []MessageData) (string, error) {
	// Собираем все сообщения в единый текст
	messageTexts := ""
	for _, msg := range messages {
		messageTexts += fmt.Sprintf("%s: %s\n", msg.SenderUserName, msg.Text)
	}

	// Формируем запрос к API
	request := OllamaRequest{
		//Model:  "llama3.1",
		Model:  "codestral",
		Prompt: fmt.Sprintf("Отвечай на русском языке. Вкратце перескажи содержимое текста ниже и выдели основные темы обсуждения:\n\n%s", messageTexts),
		Stream: false,
	}

	// Преобразуем запрос в JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// Выполняем запрос к API Ollama
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Проверяем успешность запроса
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama API returned status code %d", resp.StatusCode)
	}

	// Читаем и парсим ответ API
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	// Ожидаем, что результат будет содержать поле "response"
	summary, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format from Ollama API")
	}

	return summary, nil
}
