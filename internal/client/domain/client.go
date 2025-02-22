package domain

import "github.com/google/uuid"

// GophKeeperClientInterface - Интерфейс клиента GophKeeper
type GophKeeperClientInterface interface {
	// Login - Вход по логину и паролю, возвращает токен авторизации
	Login(login, password string) (string, error)
	// Register - Регистрация по логину и паролю, возвращает токен авторизации
	Register(login, password string) (string, error)
	// GetCerts - Возвращает публичный ключ для валидации и парсинга JWT
	GetCerts() ([]byte, error)
	// CreateText - Создает текст, возвращает идентификатор ресурса от сервера
	CreateText(session Session, content string) (uuid.UUID, error)
	// UpdateText - Обновляет существующий текст
	UpdateText(session Session, text Text) error
	// GetAllTexts - Получает все расшифрованные тексты пользователя
	GetAllTexts(session Session) ([]Text, error)
	// CreateBinary - Создает бинарные данные, возвращает идентификатор ресурса от сервера
	CreateBinary(session Session, content []byte) (uuid.UUID, error)
	// UpdateBinary - Обновляет существующие бинарные данные
	UpdateBinary(session Session, bin Binary) error
	// GetAllBinaries - Получает все расшифрованные бинарные данные пользователя
	GetAllBinaries(session Session) ([]Binary, error)
	// CreateCredentials - Создает пару логин и пароль, возвращает идентификатор ресурса от сервера
	CreateCredentials(session Session, name, login, password, meta string) (uuid.UUID, error)
	// UpdateCredentials - Обновляет существующую пару логина и пароля
	UpdateCredentials(session Session, cred *Credentials) error
	// GetAllCredentials - Получает все расшифрованные логины и пароли пользователя
	GetAllCredentials(session Session) ([]Credentials, error)
	// CreateBankCard - Создает банковскую карту, возвращает идентификатор ресурса от сервера
	CreateBankCard(session Session, number, validThru, cvv, cardHolder, meta string) (uuid.UUID, error)
	// UpdateBankCard - Обновляет существующую банковскую карту
	UpdateBankCard(session Session, card *BankCard) error
	// GetAllBankCards - Получает все расшифрованные банковские карты пользователя
	GetAllBankCards(session Session) ([]BankCard, error)
	// GetAll - Получает все расшифрованные данные пользователя
	GetAll(session Session) ([]Text, []BankCard, []Binary, []Credentials, error)
}
