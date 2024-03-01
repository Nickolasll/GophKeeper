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
	// Get - Возвращает текстовые данные по идентификатору данных и пользователя, если они существуют
	Get(textID uuid.UUID, userID uuid.UUID) (*Text, error)
	// FindByUserID - Возвращает список текстовых данных, принадлежащих пользователю
	FindByUserID(userID uuid.UUID) ([]Text, error)
}
