// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// UpdateBinary - Сценарий обновления существующих бинарных данных
type UpdateBinary struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// BinaryRepository - Реализация интерфейса BinaryRepositoryInterface
	BinaryRepository domain.BinaryRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u UpdateBinary) Do(session domain.Session, binID uuid.UUID, content []byte) error {
	bin, err := u.BinaryRepository.Get(session.UserID, binID)
	if err != nil {
		return err
	}

	bin.Content = content

	if err := u.Client.UpdateBinary(session, bin); err != nil {
		return err
	}

	if err := u.BinaryRepository.Update(session.UserID, bin); err != nil {
		return err
	}

	return nil
}
