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

// Credentials - Сущность типа хранимой информации "Логин и пароль"
type Credentials struct {
	// ID - Уникальный идентификатор "Текстовых данных"
	ID uuid.UUID
	// UserID - Ссылка на пользователя
	UserID uuid.UUID
	// Name - Зашифрованное наименование
	Name []byte
	// Login - Зашифрованный логин
	Login []byte
	// Password - Зашифрованный пароль
	Password []byte
	// Meta - Зашифрованные произвольные текстовые метаданные
	Meta []byte
}

// BankCard - Сущность типа хранимой информации "Банковкая карта"
type BankCard struct {
	// ID - Уникальный идентификатор "Банковской карты"
	ID uuid.UUID
	// UserID - Ссылка на пользователя
	UserID uuid.UUID
	// Number - Зашифрованный номер карты
	Number []byte
	// ValidThru - Зашифрованный срок действия карты
	ValidThru []byte
	// CVV - Зашифрованный CVV код карты
	CVV []byte
	// CardHolder - Зашифрованные имя и фамилия держателя карты
	CardHolder []byte
	// Meta - Зашифрованные произвольные текстовые метаданные
	Meta []byte
}
