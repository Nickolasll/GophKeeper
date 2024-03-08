package tests

import (
	"github.com/google/uuid"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// FakeHTTPClient - Фейковый клиент GophKeeper для unit тестов
type FakeHTTPClient struct {
	Certs []byte
	// Response - Возвращаемый ответ любого типа
	Response any
	// Err - Возвращаемая ошибка
	Err error
}

// Login - Вход по логину и паролю, возвращает токен авторизации
func (c FakeHTTPClient) Login(_, _ string) (string, error) {
	if c.Err != nil {
		return "", c.Err
	}

	return c.Response.(string), nil
}

// Register - Регистрация по логину и паролю, возвращает токен авторизации
func (c FakeHTTPClient) Register(_, _ string) (string, error) {
	if c.Err != nil {
		return "", c.Err
	}

	return c.Response.(string), nil
}

// CreateText - Создает текст, возвращает идентификатор ресурса от сервера
func (c FakeHTTPClient) CreateText(_ domain.Session, _ string) (uuid.UUID, error) {
	if c.Err != nil {
		return uuid.New(), c.Err
	}

	return c.Response.(uuid.UUID), nil
}

// UpdateText - Обновляет существующий текст
func (c FakeHTTPClient) UpdateText(_ domain.Session, _ domain.Text) error {
	return c.Err
}

// GetCerts - Возвращает публичный ключ для валидации и парсинга JWT
func (c FakeHTTPClient) GetCerts() ([]byte, error) {
	return c.Certs, nil
}

// CreateBinary - Создает бинарные данные, возвращает идентификатор ресурса от сервера
func (c FakeHTTPClient) CreateBinary(_ domain.Session, _ []byte) (uuid.UUID, error) {
	if c.Err != nil {
		return uuid.New(), c.Err
	}

	return c.Response.(uuid.UUID), nil
}

// UpdateText - Обновляет существующие бинарные данные
func (c FakeHTTPClient) UpdateBinary(_ domain.Session, _ domain.Binary) error {
	return c.Err
}

// CreateCredentials - Создает пару логин и пароль, возвращает идентификатор ресурса от сервера
func (c FakeHTTPClient) CreateCredentials(_ domain.Session, _, _, _ string) (uuid.UUID, error) {
	if c.Err != nil {
		return uuid.New(), c.Err
	}

	return c.Response.(uuid.UUID), nil
}

// UpdateCredentials - Обновляет существующий логин и пароль
func (c FakeHTTPClient) UpdateCredentials(_ domain.Session, _ domain.Credentials) error {
	return c.Err
}

// CreateCredentials - Создает новую банковскую карту, возвращает идентификатор ресурса от сервера
func (c FakeHTTPClient) CreateBankCard(_ domain.Session, _, _, _, _ string) (uuid.UUID, error) {
	if c.Err != nil {
		return uuid.New(), c.Err
	}

	return c.Response.(uuid.UUID), nil
}

// UpdateBankCard - Обновляет существующую банковскую карту
func (c FakeHTTPClient) UpdateBankCard(_ domain.Session, _ *domain.BankCard) error {
	return c.Err
}

// GetAllTexts - Получает все расшифрованные тексты пользователя
func (c FakeHTTPClient) GetAllTexts(_ domain.Session) ([]domain.Text, error) {
	if c.Err != nil {
		return []domain.Text{}, c.Err
	}

	return c.Response.([]domain.Text), nil
}

// GetAllBinaries - Получает все расшифрованные бинарные данные пользователя
func (c FakeHTTPClient) GetAllBinaries(_ domain.Session) ([]domain.Binary, error) {
	if c.Err != nil {
		return []domain.Binary{}, c.Err
	}

	return c.Response.([]domain.Binary), nil
}
