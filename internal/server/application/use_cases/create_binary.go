// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// CreateBinary - Сценарий использования для создания зашифрованных бинарных данных
type CreateBinary struct {
	// BinaryRepository - Интерфейс репозитория для сохранения бинарных данных
	BinaryRepository domain.BinaryRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования, возвращает идентификатор ресурса
func (u *CreateBinary) Do(userID uuid.UUID, content []byte) (uuid.UUID, error) {
	binID := uuid.New()
	encryptedContent, err := u.Crypto.Encrypt(content)
	if err != nil {
		return binID, err
	}
	text := domain.Binary{
		ID:      binID,
		UserID:  userID,
		Content: encryptedContent,
	}
	err = u.BinaryRepository.Create(text)

	return binID, err
}
