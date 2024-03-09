// Package domain содержит в себе описание сущностей и интерфейсов.
package domain

import (
	"github.com/google/uuid"
)

// UserRepositoryInterface - Интерфейс репозитория пользователя
type UserRepositoryInterface interface {
	// Create - Сохраняет нового пользователя
	Create(user User) error
	// GetByLogin - Возвращает пользователя по логину, если он существует
	GetByLogin(login string) (*User, error)
}

// TextRepositoryInterface - Интерфейс репозитория для произвольных текстовых данных
type TextRepositoryInterface interface {
	// Create - Сохраняет новые текстовые данные
	Create(text Text) error
	// Update - Сохраняет существующие текстовые данные
	Update(text Text) error
	// Get - Возвращает текстовые данные по идентификатору пользователя и данных, если они существуют
	Get(userID uuid.UUID, textID uuid.UUID) (*Text, error)
	// GetAll - Возвращает список текстовых данных, принадлежащих пользователю
	GetAll(userID uuid.UUID) ([]Text, error)
}

// BinaryRepositoryInterface - Интерфейс репозитория для произвольных бинарных данных
type BinaryRepositoryInterface interface {
	// Create - Сохраняет новые бинарные данные
	Create(bin Binary) error
	// Update - Сохраняет существующие бинарные данные
	Update(bin Binary) error
	// Get - Возвращает бинарные данные по идентификатору пользователя и данных, если они существуют
	Get(userID uuid.UUID, binID uuid.UUID) (*Binary, error)
	// GetAll - Возвращает список бинарных данных, принадлежащих пользователю
	GetAll(userID uuid.UUID) ([]Binary, error)
}

// CredentialsRepositoryInterface - Интерфейс репозитория для логинов и паролей
type CredentialsRepositoryInterface interface {
	// Create - Сохраняет новую пару логина и пароля
	Create(cred *Credentials) error
	// Update - Сохраняет существующую пару логина и пароля
	Update(cred *Credentials) error
	// Get - Возвращает пару логин и пароль по идентификатору пользователя и данных, если они существуют
	Get(userID uuid.UUID, credID uuid.UUID) (*Credentials, error)
	// GetAll - Возвращает список логинов и паролей, принадлежащих пользователю
	GetAll(userID uuid.UUID) ([]Credentials, error)
}

// BankCardRepositoryInterface - Интерфейс репозитория для банковских карт
type BankCardRepositoryInterface interface {
	// Create - Сохраняет новую банковскую карту
	Create(card *BankCard) error
	// Update - Сохраняет существующую банковскую карту
	Update(card *BankCard) error
	// Get - Возвращает банковскую карту по идентификатору пользователя и данных, если они существуют
	Get(userID uuid.UUID, cardID uuid.UUID) (*BankCard, error)
	// GetAll - Возвращает список банковских карт, принадлежащих пользователю
	GetAll(userID uuid.UUID) ([]*BankCard, error)
}
