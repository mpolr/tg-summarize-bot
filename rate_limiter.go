package main

import "time"

// Структура для отслеживания количества вызовов
type RateLimiter struct {
	UserLimits    map[int64]time.Time
	ChannelLimits map[int64]time.Time
	MaxPerMinute  int
	MaxPerHour    int
	MaxPerDay     int
}

// Создаем новый лимитер
func NewRateLimiter(maxPerMinute, maxPerHour, maxPerDay int) *RateLimiter {
	return &RateLimiter{
		UserLimits:    make(map[int64]time.Time),
		ChannelLimits: make(map[int64]time.Time),
		MaxPerMinute:  maxPerMinute,
		MaxPerHour:    maxPerHour,
		MaxPerDay:     maxPerDay,
	}
}

// Проверяем, может ли пользователь вызвать бота
func (r *RateLimiter) CanCall(userID, chatID int64) bool {
	// Здесь должна быть логика для проверки лимитов
	// Например, сколько времени прошло с последнего вызова
	return true
}
