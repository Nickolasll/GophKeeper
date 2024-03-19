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

// Credentials - Сущность типа хранимой информации "Логин и пароль"
type Credentials struct {
	// ID - Уникальный идентификатор "логина и пароля"
	ID uuid.UUID
	// Name - Наименование ресурса
	Name string
	// Login - Логин
	Login string
	// Password - пароль
	Password string
	// Meta - Зашифрованные произвольные текстовые метаданные
	Meta string
}

// BankCard - Сущность типа хранимой информации "Банковкая карта"
type BankCard struct {
	// ID - Уникальный идентификатор "Банковской карты"
	ID uuid.UUID
	// Number - Зашифрованный номер карты
	Number string
	// ValidThru - Зашифрованный срок действия карты
	ValidThru string
	// CVV - Зашифрованный CVV код карты
	CVV string
	// CardHolder - Зашифрованные имя и фамилия держателя карты
	CardHolder string
	// Meta - Зашифрованные произвольные текстовые метаданные
	Meta string
}
