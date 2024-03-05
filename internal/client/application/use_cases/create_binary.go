// Package usecases содержит имплементацию бизнес логики приложения
package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// CreateBinary - Сценарий создания новых произвольных бинарных данных
type CreateBinary struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// BinaryRepository - Реализация интерфейса BinaryRepositoryInterface
	BinaryRepository domain.BinaryRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u CreateBinary) Do(session domain.Session, content []byte) error {
	textID, err := u.Client.CreateBinary(session, content)
	if err != nil {
		return err
	}

	bin := domain.Binary{
		ID:      textID,
		Content: content,
	}

	if err := u.BinaryRepository.Create(session.UserID, bin); err != nil {
		return err
	}

	return nil
}
