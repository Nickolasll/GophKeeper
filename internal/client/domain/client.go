package domain

import "github.com/google/uuid"

// GophKeeperClientInterface - Интерфейс клиента GophKeeper
type GophKeeperClientInterface interface {
	// Login - Вход по логину и паролю, возвращает токен авторизации
	Login(login, password string) (string, error)
	// Register - Регистрация по логину и паролю, возвращает токен авторизации
	Register(login, password string) (string, error)
	// CreateText - Создает текст, возвращает идентификатор ресурса от сервера
	CreateText(session Session, content string) (uuid.UUID, error)
	// UpdateText - Обновляет существующий текст
	UpdateText(session Session, text Text) error
	// GetCerts - Возвращает публичный ключ для валидации и парсинга JWT
	GetCerts() ([]byte, error)
	// CreateBinary - Создает бинарные данные, возвращает идентификатор ресурса от сервера
	CreateBinary(session Session, content []byte) (uuid.UUID, error)
	// UpdateBinary - Обновляет существующие бинарные данные
	UpdateBinary(session Session, bin Binary) error
	// CreateCredentials - Создает пару логин и пароль, возвращает идентификатор ресурса от сервера
	CreateCredentials(session Session, name, login, password string) (uuid.UUID, error)
	// UpdateCredentials - Обновляет существующую пару логина и пароля
	UpdateCredentials(session Session, cred Credentials) error
	// CreateBankCard - Создает банковскую карту, возвращает идентификатор ресурса от сервера
	CreateBankCard(session Session, number, validThru, cvv, cardHolder string) (uuid.UUID, error)
	// UpdateBankCard - Обновляет существующую банковскую карту
	UpdateBankCard(session Session, card *BankCard) error
}
