// Package domain содержит сущности и интерфейсы к репозиториям и клиенту
package domain

import (
	"github.com/google/uuid"
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
	Create(userID uuid.UUID, text Text) error
	// Update - Сохраняет существующие текстовые данные
	Update(userID uuid.UUID, text Text) error
	// Get - Возвращает текстовые данные по идентификатору данных и пользователя, если они существуют
	Get(userID, textID uuid.UUID) (Text, error)
	// GetAll - возвращает все текстовые данные для пользователя
	GetAll(userID uuid.UUID) ([]Text, error)
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

// BinaryRepositoryInterface - Интерфейс репозитория для произвольных бинарных данных
type BinaryRepositoryInterface interface {
	// Create - Сохраняет новые бинарные данные
	Create(userID uuid.UUID, bin Binary) error
	// Update - Сохраняет существующие бинарные данные
	Update(userID uuid.UUID, bin Binary) error
	// Get - Возвращает бинарные данные по идентификатору данных и пользователя, если они существуют
	Get(userID, binID uuid.UUID) (Binary, error)
	// GetAll - возвращает все бинарные данные для пользователя
	GetAll(userID uuid.UUID) ([]Binary, error)
}
