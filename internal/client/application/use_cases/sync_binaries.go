package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// SyncBinary - Сценарий синхронизации бинарных данных
type SyncBinary struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// BinaryRepository - Реализация интерфейса BinaryRepositoryInterface
	BinaryRepository domain.BinaryRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u SyncBinary) Do(session domain.Session) error {
	bins, err := u.Client.GetAllBinaries(session)
	if err != nil {
		return err
	}

	if err := u.BinaryRepository.ReplaceAll(session.UserID, bins); err != nil {
		return err
	}

	return nil
}
