package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// SyncCredentials - Сценарий синхронизации логинов и паролей
type SyncCredentials struct {
	// Client - Реализация интерфейса GophKeeperClient
	Client domain.GophKeeperClientInterface
	// CredentialsRepository - Реализация интерфейса CredentialsRepositoryInterface
	CredentialsRepository domain.CredentialsRepositoryInterface
	// Log - логгер
	Log *logrus.Logger
}

// Do - Вызов логики сценария использования
func (u SyncCredentials) Do(session domain.Session) error {
	creds, err := u.Client.GetAllCredentials(session)
	if err != nil {
		return err
	}

	if err := u.CredentialsRepository.ReplaceAll(session.UserID, creds); err != nil {
		return err
	}

	return nil
}
