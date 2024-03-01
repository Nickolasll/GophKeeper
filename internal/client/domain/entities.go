// Package domain содержит сущности и интерфейсы к репозиториям и клиенту
package domain

// Session - Сущность сессии
type Session struct {
	// ID - Уникальный идентификатор пользователя
	UserID string
	// Token - Токен авторизации
	Token string
}

// Text - Сущность типа хранимой информации "Произвольный текст"
type Text struct {
	// ID - Уникальный идентификатор "Текстовых данных"
	ID string
	// Content - Текст
	Content string
}
