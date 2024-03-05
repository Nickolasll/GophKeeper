// Package domain содержит сущности и интерфейсы к репозиториям и клиенту
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
}
