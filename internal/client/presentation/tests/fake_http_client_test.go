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
func (c FakeHTTPClient) CreateText(_ domain.Session, _ string) (string, error) {
	if c.Err != nil {
		return uuid.NewString(), c.Err
	}

	return c.Response.(string), nil
}

// UpdateText - Обновляет существующий текст
func (c FakeHTTPClient) UpdateText(_ domain.Session, _ domain.Text) error {
	return c.Err
}

// GetCerts - Возвращает публичный ключ для валидации и парсинга JWT
func (c FakeHTTPClient) GetCerts() ([]byte, error) {
	return c.Certs, nil
}

// CreateText - Создает текст, возвращает идентификатор ресурса от сервера
func (c FakeHTTPClient) CreateBinary(_ domain.Session, _ []byte) (string, error) {
	if c.Err != nil {
		return uuid.NewString(), c.Err
	}

	return c.Response.(string), nil
}

// UpdateText - Обновляет существующий текст
func (c FakeHTTPClient) UpdateBinary(_ domain.Session, _ domain.Binary) error {
	return c.Err
}
