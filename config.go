package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	OllamaAPIURL     string `yaml:"ollama_api_url"`
	TelegramBotToken string `yaml:"telegram_bot_token"`
}

// Функция загрузки конфигурации из файла
func loadConfig(filename string) (*Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
