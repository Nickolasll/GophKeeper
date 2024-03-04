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
