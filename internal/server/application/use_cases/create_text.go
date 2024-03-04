// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// CreateText - Сценарий использования для создания зашифрованных текстовых данных
type CreateText struct {
	// TextRepository - Интерфейс репозитория для сохранения текстовых данных
	TextRepository domain.TextRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto domain.CryptoServiceInterface
}

// Execute - Вызов исполнения сценария использования, возвращает идентификатор ресурса
func (u CreateText) Execute(userID uuid.UUID, content string) (uuid.UUID, error) {
	textID := uuid.New()
	encryptedContent, err := u.Crypto.Encrypt([]byte(content))
	if err != nil {
		return textID, err
	}
	text := domain.Text{
		ID:      textID,
		UserID:  userID,
		Content: encryptedContent,
	}
	err = u.TextRepository.Create(text)

	return textID, err
}
