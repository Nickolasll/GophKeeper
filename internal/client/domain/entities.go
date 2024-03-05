// Package domain содержит сущности и интерфейсы к репозиториям и клиенту
package domain

import "github.com/google/uuid"

// Session - Сущность сессии
type Session struct {
	// ID - Уникальный идентификатор пользователя
	UserID uuid.UUID
	// Token - Токен авторизации
	Token string
}

// Text - Сущность типа хранимой информации "Произвольный текст"
type Text struct {
	// ID - Уникальный идентификатор "Текстовых данных"
	ID uuid.UUID
	// Content - Текст
	Content string
}

// Binary - Сущность типа хранимой информации "Произвольные бинарные данные"
type Binary struct {
	// ID - Уникальный идентификатор "Бинарных данных данных"
	ID uuid.UUID
	// Content - Бинарные данные
	Content []byte
}
