// Package domain содержит сущности и интерфейсы к репозиториям и клиенту
package domain

import (
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// SessionRepositoryInterface - Интерфейс репозитория сессий
type SessionRepositoryInterface interface {
	// Save - Сохраняет новую сессию
	Save(session Session) error
	// Get - Возвращает последнюю сессию, если она существует
	Get() (*Session, error)
	// Delete - Удаляет существуюущую сессию
	Delete() error
}

// TextRepositoryInterface - Интерфейс репозитория для произвольных текстовых данных
type TextRepositoryInterface interface {
	// Create - Сохраняет новые текстовые данные
	Create(userID string, text Text) error
	// Update - Сохраняет существующие текстовые данные
	Update(userID string, text Text) error
	// Get - Возвращает текстовые данные по идентификатору данных и пользователя, если они существуют
	Get(userID string, textID string) (Text, error)
	// GetAll - возвращает все текстовые данные для пользователя
	GetAll(userID string) ([]Text, error)
}

// JWKRepositoryInterface - Интерфейс хранилища публичного ключа
type JWKRepositoryInterface interface {
	// Save - Сохраняет публичный ключ
	Save(key jwk.Key) error
	// Get - Возвращает публичный ключ, если он существует
	Get() (jwk.Key, error)
	// Delete - Удаляет существующий ключ
	Delete() error
}
