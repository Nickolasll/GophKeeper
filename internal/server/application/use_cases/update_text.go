// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"

	"github.com/Nickolasll/goph-keeper/internal/server/application/services"
	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// UpdateText - Сценарий использования для обновления существующих зашифрованных текстовых данных
type UpdateText struct {
	// TextRepository - Интерфейс репозитория для сохранения текстовых данны
	TextRepository domain.TextRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto services.CryptoService
}

// Execute - Вызов исполнения сценария использования
func (u UpdateText) Execute(userID, id uuid.UUID, content string) error {
	text, err := u.TextRepository.Get(id, userID)
	if err != nil {
		return err
	}
	if text == nil {
		return domain.ErrEntityNotFound
	}

	encryptedContent, err := u.Crypto.Encrypt(content)
	if err != nil {
		return err
	}
	text.Content = encryptedContent
	err = u.TextRepository.Update(*text)

	return err
}
