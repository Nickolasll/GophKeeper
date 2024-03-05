// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// UpdateBinary - Сценарий использования для обновления существующих зашифрованных бинарных данных
type UpdateBinary struct {
	// TextRepository - Интерфейс репозитория для сохранения бинарных данных
	BinaryRepository domain.BinaryRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto domain.CryptoServiceInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов исполнения сценария использования
func (u UpdateBinary) Do(userID, id uuid.UUID, content []byte) error {
	bin, err := u.BinaryRepository.Get(userID, id)
	if err != nil {
		return err
	}
	if bin == nil {
		return domain.ErrEntityNotFound
	}

	encryptedContent, err := u.Crypto.Encrypt(content)
	if err != nil {
		return err
	}
	bin.Content = encryptedContent
	err = u.BinaryRepository.Update(*bin)

	return err
}
