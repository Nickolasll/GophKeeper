// Package domain содержит в себе описание сущностей и интерфейсов.
package domain

import (
	"github.com/google/uuid"
)

// User - Сущность пользователя, используется для авторизации
type User struct {
	// ID - Уникальный идентификатор пользователя
	ID uuid.UUID
	// Login - Логин пользователя
	Login string
	// Password - Хэш пароля пользователя
	Password string
}

// Text - Сущность типа хранимой информации "Произвольный текст"
type Text struct {
	// ID - Уникальный идентификатор "Текстовых данных"
	ID uuid.UUID
	// UserID - Ссылка на пользователя
	UserID uuid.UUID
	// Content - Зашифрованный текст
	Content []byte
}

// Binary - Сущность типа хранимой информации "Произвольные бинарные данные"
type Binary struct {
	// ID - Уникальный идентификатор "Текстовых данных"
	ID uuid.UUID
	// UserID - Ссылка на пользователя
	UserID uuid.UUID
	// Content - Зашифрованные бинарные данные
	Content []byte
}
