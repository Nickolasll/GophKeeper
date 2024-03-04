// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"

	"github.com/Nickolasll/goph-keeper/internal/server/application/services"
	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// CreateBinary - Сценарий использования для создания зашифрованных бинарных данных
type CreateBinary struct {
	// TextRepository - Интерфейс репозитория для сохранения бинарных данных
	BinaryRepository domain.BinaryRepositoryInterface
	// Crypto - Сервис для шифрования данных
	Crypto services.CryptoService
}

// Execute - Вызов исполнения сценария использования, возвращает идентификатор ресурса
func (u *CreateBinary) Execute(userID uuid.UUID, content []byte) (uuid.UUID, error) {
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
