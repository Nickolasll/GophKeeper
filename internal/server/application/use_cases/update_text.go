package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// UpdateText - Сценарий использования для обновления существующих зашифрованных текстовых данных
type UpdateText struct {
	// TextRepository - Интерфейс репозитория для сохранения текстовых данных
	TextRepository domain.TextRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования
func (u UpdateText) Do(userID, id uuid.UUID, content string) error {
	text, err := u.TextRepository.Get(userID, id)
	if err != nil {
		return err
	}
	if text == nil {
		return domain.ErrEntityNotFound
	}

	encryptedContent, err := u.Crypto.Encrypt([]byte(content))
	if err != nil {
		return err
	}
	text.Content = encryptedContent
	err = u.TextRepository.Update(*text)

	return err
}
