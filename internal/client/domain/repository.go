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
	// GetAll - Возвращает все текстовые данные для пользователя
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
	// GetAll - Возвращает все бинарные данные для пользователя
	GetAll(userID uuid.UUID) ([]Binary, error)
}

// CredentialsRepositoryInterface - Интерфейс репозитория для логинов и паролей
type CredentialsRepositoryInterface interface {
	// Create - Сохраняет новую пару логина и пароля
	Create(userID uuid.UUID, cred Credentials) error
	// Update - Сохраняет существующую пару логина и пароля
	Update(userID uuid.UUID, cred Credentials) error
	// Get - Возвращает пару логин и парль по идентификатору данных и пользователя, если они существуют
	Get(userID, credID uuid.UUID) (Credentials, error)
	// GetAll - Возвращает все логины и пароли для пользователя
	GetAll(userID uuid.UUID) ([]Credentials, error)
}

// BankCardRepositoryInterface - Интерфейс репозитория для банковских карт
type BankCardRepositoryInterface interface {
	// Create - Сохраняет новую банковскую карту
	Create(userID uuid.UUID, card *BankCard) error
	// Update - Сохраняет существующую банковскую карту
	Update(userID uuid.UUID, card *BankCard) error
	// Get - Возвращает банковскую карту по идентификатору данных и пользователя, если они существуют
	Get(userID, cardID uuid.UUID) (BankCard, error)
	// GetAll - Возвращает все банковские карты для пользователя
	GetAll(userID uuid.UUID) ([]BankCard, error)
}
